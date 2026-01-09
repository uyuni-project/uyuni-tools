// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CreatePersistentVolumeClaims creates all the PVCs described by the mounts.
func CreatePersistentVolumeClaims(
	namespace string,
	mounts []types.VolumeMount,
) error {
	pvcs := GetPersistentVolumeClaims(
		namespace,
		"",
		core.ReadWriteOnce,
		false,
		GetLabels(ServerApp, ""),
		mounts,
	)

	for _, pvc := range pvcs {
		if !hasPersistentVolumeClaim(pvc.Namespace, pvc.Name) {
			if err := Apply(
				[]*core.PersistentVolumeClaim{pvc},
				fmt.Sprintf(L("failed to create %s persistent volume claim"), pvc.Name),
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func hasPersistentVolumeClaim(namespace string, name string) bool {
	out, err := runCmdOutput(zerolog.DebugLevel, "kubectl", "get", "pvc", "-n", namespace, name, "-o", "name")
	return err == nil && strings.TrimSpace(string(out)) != ""
}

// CreatePersistentVolumeClaimForVolume creates a PVC bound to a specific Volume.
func CreatePersistentVolumeClaimForVolume(
	namespace string,
	volumeName string,
) error {
	// Get the PV Storage class and claimRef
	out, err := utils.RunCmdOutput(zerolog.DebugLevel,
		"kubectl", "get", "pv", volumeName, "-n", namespace, "-o", "json",
	)
	if err != nil {
		return err
	}
	var pv core.PersistentVolume
	if err := json.Unmarshal(out, &pv); err != nil {
		return utils.Errorf(err, L("failed to parse pv data"))
	}

	// Ensure the claimRef of the volume is for our PVC
	if pv.Spec.ClaimRef == nil || pv.Spec.ClaimRef.Name != volumeName && pv.Spec.ClaimRef.Namespace != namespace {
		return fmt.Errorf(L("the %[1]s volume has to reference the %[1]s claim in %[2]s namespace"), volumeName, namespace)
	}

	// Create the PVC object
	pvc := newPersistentVolumeClaim(
		namespace, volumeName, pv.Spec.StorageClassName,
		pv.Spec.Capacity.Storage().String(), pv.Spec.AccessModes, false,
	)
	pvc.Spec.VolumeName = volumeName

	return Apply([]runtime.Object{&pvc}, L("failed to run the persistent volume claims"))
}

// GetPersistentVolumeClaims creates the PVC objects matching a list of volume mounts.
func GetPersistentVolumeClaims(
	namespace string,
	storageClass string,
	accessMode core.PersistentVolumeAccessMode,
	matchPvByLabel bool,
	labels map[string]string,
	mounts []types.VolumeMount,
) []*core.PersistentVolumeClaim {
	var claims []*core.PersistentVolumeClaim

	for _, mount := range mounts {
		size := mount.Size
		if size == "" {
			log.Warn().Msgf(L("no size defined for PersistentVolumeClaim %s, using 10Mi as default"), mount.Name)
			size = "10Mi"
		}
		pv := newPersistentVolumeClaim(
			namespace,
			mount.Name,
			storageClass,
			size,
			[]core.PersistentVolumeAccessMode{accessMode},
			matchPvByLabel,
		)
		pv.SetLabels(labels)
		claims = append(claims, &pv)
	}

	return claims
}

// Creates a PVC from a few common values.
func newPersistentVolumeClaim(
	namespace string,
	name string,
	storageClass string,
	size string,
	accessModes []core.PersistentVolumeAccessMode,
	matchPvByLabel bool,
) core.PersistentVolumeClaim {
	pvc := core.PersistentVolumeClaim{
		TypeMeta: v1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolumeClaim",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: core.PersistentVolumeClaimSpec{
			AccessModes: accessModes,
			Resources: core.VolumeResourceRequirements{
				Requests: core.ResourceList{"storage": resource.MustParse(size)},
			},
		},
	}

	if storageClass != "" {
		pvc.Spec.StorageClassName = &storageClass
	}

	if matchPvByLabel {
		pvc.Spec.Selector = &v1.LabelSelector{
			MatchLabels: map[string]string{"data": name},
		}
	}

	return pvc
}

func createMount(mountPath string) core.VolumeMount {
	pattern := regexp.MustCompile("[^a-zA-Z]+")
	name := strings.Trim(pattern.ReplaceAllString(mountPath, "-"), "-")
	return core.VolumeMount{
		MountPath: mountPath,
		Name:      name,
	}
}

// CreateTmpfsMount creates a temporary volume and its mount.
func CreateTmpfsMount(mountPath string, size string) (core.VolumeMount, core.Volume) {
	mount := createMount(mountPath)

	parsedSize := resource.MustParse(size)
	volume := core.Volume{
		Name: mount.Name,
		VolumeSource: core.VolumeSource{
			EmptyDir: &core.EmptyDirVolumeSource{
				Medium:    core.StorageMediumMemory,
				SizeLimit: &parsedSize,
			},
		},
	}
	return mount, volume
}

// CreateHostPathMount creates the mount and volume for a host path.
// This is not secure and tied to the availability on the node, only use when needed.
func CreateHostPathMount(
	mountPath string,
	hostPath string,
	sourceType core.HostPathType,
) (core.VolumeMount, core.Volume) {
	mount := createMount(mountPath)

	volume := core.Volume{
		Name: mount.Name,
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: hostPath,
				Type: &sourceType,
			},
		},
	}
	return mount, volume
}

// CreateSecretMount creates the volume for a secret.
func CreateSecretVolume(name string, secretName string) core.Volume {
	volume := core.Volume{
		Name: name,
		VolumeSource: core.VolumeSource{
			Secret: &core.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	}

	return volume
}

// CreateConfigVolume creates the volume for a ConfigMap.
func CreateConfigVolume(name string, configMapName string) core.Volume {
	volume := core.Volume{
		Name: name,
		VolumeSource: core.VolumeSource{
			ConfigMap: &core.ConfigMapVolumeSource{
				LocalObjectReference: core.LocalObjectReference{
					Name: configMapName,
				},
			},
		},
	}

	return volume
}

// CreateVolumes creates PVC-based volumes matching the internal volumes mounts.
func CreateVolumes(mounts []types.VolumeMount) []core.Volume {
	volumes := []core.Volume{}

	for _, mount := range mounts {
		volume := core.Volume{
			Name: mount.Name,
			VolumeSource: core.VolumeSource{
				PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
					ClaimName: mount.Name,
				},
			},
		}
		volumes = append(volumes, volume)
	}

	return volumes
}

var runCmdOutput = utils.RunCmdOutput

// HasVolume returns true if the pvcName persistent volume claim is bound.
func HasVolume(namespace string, pvcName string) bool {
	out, err := runCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "pvc", "-n", namespace, pvcName, "-o", "jsonpath={.status.phase}",
	)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "Bound"
}
