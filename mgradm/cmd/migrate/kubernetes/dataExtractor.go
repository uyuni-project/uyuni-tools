// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"errors"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"gopkg.in/yaml.v2"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// MigrationData represents the files and data extracted from the migration sync phase.
type MigrationData struct {
	CaKey      string
	CaCert     string
	Data       *utils.InspectResult
	ServerCert string
	ServerKey  string
}

func extractMigrationData(namespace string, image string, volume types.VolumeMount) (*MigrationData, error) {
	// Read the file from the volume from a container into stdout
	mounts := kubernetes.ConvertVolumeMounts([]types.VolumeMount{volume})
	volumes := kubernetes.CreateVolumes([]types.VolumeMount{volume})

	podName := "uyuni-data-extractor"

	// Use a pod here since this is a very simple task reading out a file from a volume
	pod := core.Pod{
		TypeMeta:   meta.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: meta.ObjectMeta{Name: podName, Namespace: namespace},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "extractor",
					Image:           image,
					ImagePullPolicy: core.PullIfNotPresent,
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

	tempDir, err := utils.TempDir()
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

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

	// Parse the content
	files := make(map[string]string)
	if err := yaml.Unmarshal(data, &files); err != nil {
		return nil, utils.Errorf(err, L("failed to parse data extractor pod output"))
	}

	var result MigrationData
	for file, content := range files {
		if file == "RHN-ORG-PRIVATE-SSL-KEY" {
			result.CaKey = content
		} else if file == "RHN-ORG-TRUSTED-SSL-CERT" {
			result.CaCert = content
		} else if file == "spacewalk.crt" {
			result.ServerCert = content
		} else if file == "spacewalk.key" {
			result.ServerKey = content
		} else if file == "data" {
			parsedData, err := utils.ReadInspectDataString[utils.InspectResult]([]byte(content))
			if err != nil {
				return nil, utils.Errorf(err, L("failed to parse migration data file"))
			}
			result.Data = parsedData
		}
	}

	if result.Data == nil {
		return nil, errors.New(L("found no data file after migration"))
	}

	return &result, nil
}
