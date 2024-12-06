// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"fmt"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// RunPodLogs runs a pod, waits for it to finish and returns it logs.
//
// This should be used only to run very fast tasks.
func RunPodLogs(
	namespace string,
	name string,
	image string,
	pullPolicy string,
	pullSecret string,
	volumesMounts []types.VolumeMount,
	cmd ...string,
) ([]byte, error) {
	// Read the file from the volume from a container into stdout
	mounts := ConvertVolumeMounts(volumesMounts)
	volumes := CreateVolumes(volumesMounts)

	// Use a pod here since this is a very simple task reading out a file from a volume
	pod := core.Pod{
		TypeMeta: meta.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"app": name},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            name,
					Image:           image,
					ImagePullPolicy: GetPullPolicy(pullPolicy),
					Command:         cmd,
					VolumeMounts:    mounts,
				},
			},
			Volumes:       volumes,
			RestartPolicy: core.RestartPolicyNever,
		},
	}

	if pullSecret != "" {
		pod.Spec.ImagePullSecrets = []core.LocalObjectReference{{Name: pullSecret}}
	}

	tempDir, cleaner, err := utils.TempDir()
	if err != nil {
		return nil, err
	}
	defer cleaner()

	// Run the pod
	podPath := path.Join(tempDir, "pod.yaml")
	if err := YamlFile([]runtime.Object{&pod}, podPath); err != nil {
		return nil, err
	}

	if err := utils.RunCmd("kubectl", "apply", "-f", podPath); err != nil {
		return nil, utils.Errorf(err, L("failed to run the %s pod"), name)
	}
	if err := Apply(
		[]runtime.Object{&pod}, fmt.Sprintf(L("failed to run the %s pod"), name),
	); err != nil {
		return nil, err
	}

	if err := WaitForPod(namespace, name, 60); err != nil {
		return nil, err
	}

	data, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "logs", "-n", namespace, name)
	if err != nil {
		return nil, utils.Errorf(err, L("failed to get the %s pod logs"), name)
	}

	defer func() {
		if err := DeletePod(namespace, name, "-lapp="+name); err != nil {
			log.Err(err).Msgf(L("failed to delete the %s pod"), name)
		}
	}()

	return data, nil
}
