// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func inspectServer(namespace string, image string, pullPolicy string) (*utils.ServerInspectData, error) {
	podName := "uyuni-image-inspector"

	tempDir, err := utils.TempDir()
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)
	inspector := utils.NewServerInspector(tempDir)
	// We need the inspector to write to the pod's logs instead of a file
	inspector.DataPath = "/dev/stdout"
	if err := inspector.GenerateScript(); err != nil {
		return nil, err
	}

	// Mount the postgresql and config volumes to extract data from the running instance too
	mounts := kubernetes.ConvertVolumeMounts(utils.PgsqlRequiredVolumeMounts)
	volumes := kubernetes.CreateVolumes(utils.PgsqlRequiredVolumeMounts)

	// Use a pod here since this is a very simple task reading out a file from a volume
	pod := core.Pod{
		TypeMeta:   meta.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: meta.ObjectMeta{Name: podName, Namespace: namespace},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "inspector",
					Image:           image,
					ImagePullPolicy: kubernetes.GetPullPolicy(pullPolicy),
					Command: []string{
						"sh", "-c",
						"for f in /var/lib/uyuni-tools/*; do echo \"`basename $f`: |2\"; cat $f | sed 's/^/  /'; done",
					},
					VolumeMounts: mounts,
				},
			},
			Volumes:       volumes,
			RestartPolicy: core.RestartPolicyNever,
		},
	}

	// Run the pod
	extractorPodPath := path.Join(tempDir, "extractor-pod.yaml")
	if err := kubernetes.YamlFile([]runtime.Object{&pod}, extractorPodPath); err != nil {
		return nil, err
	}

	if err := utils.RunCmd("kubectl", "apply", "-f", extractorPodPath); err != nil {
		return nil, utils.Errorf(err, L("failed to run the migration data extractor pod"))
	}
	if err := kubernetes.Apply(
		[]runtime.Object{&pod}, L("failed to run the migration data extractor pod"),
	); err != nil {
		return nil, err
	}

	if err := kubernetes.WaitForPod(namespace, podName, 60); err != nil {
		return nil, err
	}

	data, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "logs", "-n", namespace, podName)
	if err != nil {
		return nil, utils.Errorf(err, L("failed to get the migration data extractor pod logs"))
	}

	defer func() {
		if err := utils.RunCmd("kubectl", "delete", "pod", "-n", namespace, podName); err != nil {
			log.Err(err).Msgf(L("failed to delete the uyuni-data-extractor pod"))
		}
	}()

	// Parse the data
	inspectedData, err := utils.ReadInspectDataString[utils.ServerInspectData]([]byte(data))
	if err != nil {
		return nil, utils.Errorf(err, L("failed to parse the inspected data"))
	}
	return inspectedData, nil
}
