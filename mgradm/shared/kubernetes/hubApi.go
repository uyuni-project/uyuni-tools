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

const (
	// HubAPIDeployName is the deployment name of the Hub API.
	HubAPIDeployName  = "uyuni-hub-api"
	hubAPIServiceName = "hub-api"
)

// InstallHubAPI installs the Hub API deployment and service.
func InstallHubAPI(namespace string, image string, pullPolicy string, pullSecret string) error {
	if err := startHubAPIDeployment(namespace, image, pullPolicy, pullSecret); err != nil {
		return err
	}

	if err := createHubAPIService(namespace); err != nil {
		return err
	}

	// TODO Do we want an ingress to use port 80 / 443 from the outside too?
	// This would have an impact on the user's scripts.
	return nil
}

func startHubAPIDeployment(namespace string, image string, pullPolicy string, pullSecret string) error {
	deploy := getHubAPIDeployment(namespace, image, pullPolicy, pullSecret)
	return kubernetes.Apply([]runtime.Object{deploy}, L("failed to create the hub API deployment"))
}

func getHubAPIDeployment(namespace string, image string, pullPolicy string, pullSecret string) *apps.Deployment {
	var replicas int32 = 1

	deploy := &apps.Deployment{
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      HubAPIDeployName,
			Namespace: namespace,
			Labels:    kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.HubAPIComponent),
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Selector: &meta.LabelSelector{
				MatchLabels: kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.HubAPIComponent),
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.HubAPIComponent),
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
								{Name: "HUB_API_URL", Value: fmt.Sprintf("http://%s/rpc/api", utils.WebServiceName)},
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

	if pullSecret != "" {
		deploy.Spec.Template.Spec.ImagePullSecrets = []core.LocalObjectReference{{Name: pullSecret}}
	}
	return deploy
}

func createHubAPIService(namespace string) error {
	svc := getService(namespace, kubernetes.ServerApp, kubernetes.HubAPIComponent, hubAPIServiceName, core.ProtocolTCP,
		utils.NewPortMap(utils.HubAPIServiceName, "api", 2830, 2830),
	)
	return kubernetes.Apply([]runtime.Object{svc}, L("failed to create the hub API service"))
}
