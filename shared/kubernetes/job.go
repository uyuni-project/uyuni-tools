// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"strings"
	"time"

	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetScriptJob prepares the definition of a kubernetes job running a shell script from a template.
// The name is suffixed with a time stamp to avoid collisions.
func GetScriptJob(
	namespace string,
	name string,
	image string,
	pullPolicy string,
	pullSecret string,
	mounts []types.VolumeMount,
	template utils.Template,
) (*batch.Job, error) {
	var maxFailures int32

	// Convert our mounts to Kubernetes objects
	volumeMounts := ConvertVolumeMounts(mounts)
	volumes := CreateVolumes(mounts)

	// Prepare the script
	scriptBuilder := new(strings.Builder)
	if err := template.Render(scriptBuilder); err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("20060102150405")

	// Create the job object running the script wrapped as a sh command
	job := batch.Job{
		TypeMeta: meta.TypeMeta{Kind: "Job", APIVersion: "batch/v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      name + "-" + timestamp,
			Namespace: namespace,
			Labels:    GetLabels(ServerApp, ""),
		},
		Spec: batch.JobSpec{
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            "runner",
							Image:           image,
							ImagePullPolicy: GetPullPolicy(pullPolicy),
							Command:         []string{"sh", "-c", scriptBuilder.String()},
							VolumeMounts:    volumeMounts,
						},
					},
					Volumes:       volumes,
					RestartPolicy: core.RestartPolicyNever,
				},
			},
			BackoffLimit: &maxFailures,
		},
	}

	if pullSecret != "" {
		job.Spec.Template.Spec.ImagePullSecrets = []core.LocalObjectReference{{Name: pullSecret}}
	}

	return &job, nil
}
