// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"

	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ServerDeployName is the name of the server deployment.
const ServerDeployName = "uyuni"

// CreateServerDeployment creates a new deployment of the server.
func CreateServerDeployment(
	namespace string,
	serverImage string,
	pullPolicy string,
	timezone string,
	debug bool,
	mirrorPvName string,
	pullSecret string,
) error {
	if mirrorPvName != "" {
		// Create a PVC using the required mirror PV
		if err := kubernetes.CreatePersistentVolumeClaimForVolume(namespace, mirrorPvName); err != nil {
			return err
		}
	}

	serverDeploy := GetServerDeployment(
		namespace, serverImage, kubernetes.GetPullPolicy(pullPolicy), timezone, debug, mirrorPvName, pullSecret,
	)

	return kubernetes.Apply([]runtime.Object{serverDeploy}, L("failed to create the server deployment"))
}

// GetServerDeployment computes the deployment object for an Uyuni server.
func GetServerDeployment(
	namespace string,
	image string,
	pullPolicy core.PullPolicy,
	timezone string,
	debug bool,
	mirrorPvName string,
	pullSecret string,
) *apps.Deployment {
	var replicas int32 = 1

	runMount, runVolume := kubernetes.CreateTmpfsMount("/run", "256Mi")
	cgroupMount, cgroupVolume := kubernetes.CreateHostPathMount(
		"/sys/fs/cgroup", "/sys/fs/cgroup", core.HostPathDirectory,
	)

	// Compute the needed ports
	ports := utils.GetServerPorts(debug)

	template := getServerPodTemplate(image, pullPolicy, timezone, pullSecret)

	template.Spec.Volumes = append(template.Spec.Volumes, runVolume, cgroupVolume)
	template.Spec.Containers[0].Ports = kubernetes.ConvertPortMaps(ports)
	template.Spec.Containers[0].VolumeMounts = append(template.Spec.Containers[0].VolumeMounts,
		runMount, cgroupMount,
	)

	if mirrorPvName != "" {
		// Add a mount for the mirror
		template.Spec.Containers[0].VolumeMounts = append(template.Spec.Containers[0].VolumeMounts,
			core.VolumeMount{
				Name:      mirrorPvName,
				MountPath: "/mirror",
			},
		)

		// Add the environment variable for the deployment to use the mirror
		// This doesn't makes sense for migration as the setup script is not executed
		template.Spec.Containers[0].Env = append(template.Spec.Containers[0].Env,
			core.EnvVar{Name: "MIRROR_PATH", Value: "/mirror"},
		)
	}

	template.Spec.Containers[0].Lifecycle = &core.Lifecycle{
		PreStop: &core.LifecycleHandler{
			Exec: &core.ExecAction{
				Command: []string{"/bin/sh", "-c", "spacewalk-service stop && systemctl stop postgresql"},
			},
		},
	}

	template.Spec.Containers[0].ReadinessProbe = &core.Probe{
		ProbeHandler: core.ProbeHandler{
			HTTPGet: &core.HTTPGetAction{
				Port: intstr.FromInt(80),
				Path: "/rhn/manager/api/api/getVersion",
			},
		},
		PeriodSeconds:    30,
		TimeoutSeconds:   20,
		FailureThreshold: 5,
	}

	template.Spec.Containers[0].LivenessProbe = &core.Probe{
		ProbeHandler: core.ProbeHandler{
			HTTPGet: &core.HTTPGetAction{
				Port: intstr.FromInt(80),
				Path: "/rhn/manager/api/api/getVersion",
			},
		},
		InitialDelaySeconds: 60,
		PeriodSeconds:       60,
		TimeoutSeconds:      20,
		FailureThreshold:    5,
	}

	deployment := apps.Deployment{
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      ServerDeployName,
			Namespace: namespace,
			Labels:    kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.ServerComponent),
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			// As long as the container cannot scale, we need to stick to recreate strategy
			// or the new deployed pods won't be ready.
			Strategy: apps.DeploymentStrategy{Type: apps.RecreateDeploymentStrategyType},
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{kubernetes.ComponentLabel: kubernetes.ServerComponent},
			},
			Template: template,
		},
	}

	return &deployment
}

// GetServerPodTemplate computes the pod template with the init container and the minimum viable volumes and mounts.
// This is intended to be shared with the setup job.
func getServerPodTemplate(
	image string,
	pullPolicy core.PullPolicy,
	timezone string,
	pullSecret string,
) core.PodTemplateSpec {
	envs := []core.EnvVar{
		{Name: "TZ", Value: timezone},
	}

	mounts := GetServerMounts()

	// Convert our mounts to Kubernetes objects
	volumeMounts := kubernetes.ConvertVolumeMounts(mounts)

	// The init mounts are the same mounts but in /mnt just for the init container populating the volumes
	var initMounts []core.VolumeMount
	for _, mount := range volumeMounts {
		initMount := mount.DeepCopy()
		initMount.MountPath = "/mnt" + initMount.MountPath
		initMounts = append(initMounts, *initMount)
	}

	volumes := kubernetes.CreateVolumes(mounts)

	caMount := core.VolumeMount{
		Name:      "ca-cert",
		MountPath: "/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT",
		ReadOnly:  true,
		SubPath:   "ca.crt",
	}
	tlsKeyMount := core.VolumeMount{Name: "tls-key", MountPath: "/etc/pki/spacewalk-tls"}

	caVolume := kubernetes.CreateConfigVolume("ca-cert", "uyuni-ca")
	tlsKeyVolume := kubernetes.CreateSecretVolume("tls-key", "uyuni-cert")
	var keyMode int32 = 0o400
	tlsKeyVolume.VolumeSource.Secret.Items = []core.KeyToPath{
		{Key: "tls.crt", Path: "spacewalk.crt"},
		{Key: "tls.key", Path: "spacewalk.key", Mode: &keyMode},
	}

	initMounts = append(initMounts, tlsKeyMount)
	volumeMounts = append(volumeMounts, caMount, tlsKeyMount)
	volumes = append(volumes, caVolume, tlsKeyVolume)

	template := core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{
			Labels: kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.ServerComponent),
		},
		Spec: core.PodSpec{
			InitContainers: []core.Container{
				{
					Name:            "init-volumes",
					Image:           image,
					ImagePullPolicy: pullPolicy,
					Command:         []string{"sh", "-x", "-c", initScript},
					VolumeMounts:    initMounts,
				},
			},
			Containers: []core.Container{
				{
					Name:            "uyuni",
					Image:           image,
					ImagePullPolicy: pullPolicy,
					Env:             envs,
					VolumeMounts:    volumeMounts,
				},
			},
			Volumes: volumes,
		},
	}

	if pullSecret != "" {
		template.Spec.ImagePullSecrets = []core.LocalObjectReference{{Name: pullSecret}}
	}
	return template
}

const initScript = `
# Fill he empty volumes
for vol in /var/lib/cobbler \
		   /var/lib/salt \
		   /var/lib/pgsql \
		   /var/cache \
		   /var/log \
		   /srv/salt \
		   /srv/www \
		   /srv/tftpboot \
		   /srv/formula_metadata \
		   /srv/pillar \
		   /srv/susemanager \
		   /srv/spacewalk \
		   /root \
		   /etc/apache2 \
		   /etc/rhn \
		   /etc/systemd/system/multi-user.target.wants \
		   /etc/systemd/system/sockets.target.wants \
		   /etc/salt \
		   /etc/tomcat \
		   /etc/cobbler \
		   /etc/sysconfig \
		   /etc/postfix \
		   /etc/sssd \
		   /etc/pki/tls
do
	chown --reference=$vol /mnt$vol;
	chmod --reference=$vol /mnt$vol;
	if [ -z "$(ls -A /mnt$vol)" ]; then
    	cp -a $vol/. /mnt$vol;
		if [ "$vol" = "/srv/www" ]; then
            ln -s /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT /mnt$vol/RHN-ORG-TRUSTED-SSL-CERT;
		fi

		if [ "$vol" = "/etc/pki/tls" ]; then
              ln -s /etc/pki/spacewalk-tls/spacewalk.crt /mnt/etc/pki/tls/certs/spacewalk.crt;
              ln -s /etc/pki/spacewalk-tls/spacewalk.key /mnt/etc/pki/tls/private/spacewalk.key;
		fi
	fi

	if [ "$vol" = "/etc/pki/tls" ]; then
	    cp /etc/pki/spacewalk-tls/spacewalk.key /mnt/etc/pki/tls/private/pg-spacewalk.key;
	    chown postgres:postgres /mnt/etc/pki/tls/private/pg-spacewalk.key;
	fi
done
`

// GetServerMounts returns the volume mounts required for the server pod.
func GetServerMounts() []types.VolumeMount {
	// Filter out the duplicate mounts to avoid issues applying the jobs
	serverMounts := utils.ServerVolumeMounts
	mounts := []types.VolumeMount{}
	mountsSet := map[string]types.VolumeMount{}
	for _, mount := range serverMounts {
		switch mount.Name {
		// Skip mounts that are not PVCs
		case "ca-cert", "tls-key":
			continue
		}
		if _, exists := mountsSet[mount.Name]; !exists {
			mounts = append(mounts, mount)
			mountsSet[mount.Name] = mount
		}
	}

	return mounts
}

// TuneMounts adjusts the server mounts with the size and storage class passed by as parameters.
func TuneMounts(mounts []types.VolumeMount, flags *cmd_utils.VolumesFlags) []types.VolumeMount {
	tunedMounts := []types.VolumeMount{}
	for _, mount := range mounts {
		class := flags.Class
		var volumeFlags *cmd_utils.VolumeFlags
		switch mount.Name {
		case "var-pgsql":
			volumeFlags = &flags.Database
		case "var-spacewalk":
			volumeFlags = &flags.Packages
		case "var-cache":
			volumeFlags = &flags.Cache
		case "srv-www":
			volumeFlags = &flags.Www
		}
		if volumeFlags != nil {
			if volumeFlags.Class != "" {
				class = volumeFlags.Class
			}
			mount.Size = volumeFlags.Size
		}
		mount.Class = class
		tunedMounts = append(tunedMounts, mount)
	}
	return tunedMounts
}

var runCmdOutput = utils.RunCmdOutput

// getRunningServerImage extracts the main server container image from a running deployment.
func getRunningServerImage(namespace string) string {
	out, err := runCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "deploy", "-n", namespace, ServerDeployName,
		"-o", "jsonpath={.spec.template.spec.containers[0].image}",
	)
	if err != nil {
		// Errors could be that the namespace or deployment doesn't exist, just return no image.
		log.Debug().Err(err).Msg("failed to get the running server container image")
		return ""
	}
	return strings.TrimSpace(string(out))
}
