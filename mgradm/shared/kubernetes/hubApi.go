// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// HubApiDeployName is the deployment name of the Hub API.
const HubApiDeployName = "uyuni-hub-api"

const (
	hubApiAppName     = "uyuni-hub-api"
	hubApiServiceName = "hub-api"
)

// InstallHubApi installs the Hub API deployment and service.
func InstallHubApi(namespace string, image string, pullPolicy string) error {
	if err := startHubApiDeployment(namespace, image, pullPolicy); err != nil {
		return err
	}

	if err := createHubApiService(namespace); err != nil {
		return err
	}

	// TODO Do we want an ingress to use port 80 / 443 from the outside too?
	// This would have an impact on the user's scripts.
	return nil
}

func startHubApiDeployment(namespace string, image string, pullPolicy string) error {
	deploy := getHubApiDeployment(namespace, image, pullPolicy)
	return kubernetes.Apply([]runtime.Object{deploy}, L("failed to create the hub API deployment"))
}

func getHubApiDeployment(namespace string, image string, pullPolicy string) *apps.Deployment {
	var replicas int32 = 1

	return &apps.Deployment{
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      HubApiDeployName,
			Namespace: namespace,
			Labels:    map[string]string{"app": hubApiAppName},
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"app": hubApiAppName},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"app": hubApiAppName},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            "uyuni-hub-api",
							Image:           image,
							ImagePullPolicy: kubernetes.GetPullPolicy(pullPolicy),
							Ports: []core.ContainerPort{
								{
									ContainerPort: int32(2830),
								},
							},
							Env: []core.EnvVar{
								{Name: "HUB_API_URL", Value: fmt.Sprintf("http://%s/rpc/api", webServiceName)},
								{Name: "HUB_CONNECT_TIMEOUT", Value: "10"},
								{Name: "HUB_REQUEST_TIMEOUT", Value: "10"},
								{Name: "HUB_CONNECT_USING_SSL", Value: "false"},
							},
						},
					},
				},
			},
		},
	}
}

func createHubApiService(namespace string) error {
	svc := getService(namespace, hubApiServiceName, core.ProtocolTCP, utils.NewPortMap("api", 2830, 2830))
	return kubernetes.Apply([]runtime.Object{svc}, L("failed to create the hub API service"))
}
