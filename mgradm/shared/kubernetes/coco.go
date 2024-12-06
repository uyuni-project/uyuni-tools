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
	// CocoApiDeployName is the deployment name for confidential computing attestations.
	CocoDeployName = "uyuni-coco-attestation"
)

// StartCocoDeployment installs the confidential computing deployment.
func StartCocoDeployment(
	namespace string,
	image string,
	pullPolicy string,
	pullSecret string,
	replicas int,
	dbPort int,
	dbName string,
) error {
	deploy := getCocoDeployment(namespace, image, pullPolicy, pullSecret, int32(replicas), dbPort, dbName)
	return kubernetes.Apply([]runtime.Object{deploy},
		L("failed to create confidential computing attestations deployment"),
	)
}

func getCocoDeployment(
	namespace string,
	image string,
	pullPolicy string,
	pullSecret string,
	replicas int32,
	dbPort int,
	dbName string,
) *apps.Deployment {
	cnxURL := fmt.Sprintf("jdbc:postgresql://%s:%d/%s", utils.DBServiceName, dbPort, dbName)
	deploy := &apps.Deployment{
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      CocoDeployName,
			Namespace: namespace,
			Labels:    kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.CocoComponent),
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Selector: &meta.LabelSelector{
				MatchLabels: kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.CocoComponent),
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.CocoComponent),
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            "coco",
							Image:           image,
							ImagePullPolicy: kubernetes.GetPullPolicy(pullPolicy),
							Env: []core.EnvVar{
								{Name: "database_connection", Value: cnxURL},
								{Name: "database_user", ValueFrom: &core.EnvVarSource{
									SecretKeyRef: &core.SecretKeySelector{
										LocalObjectReference: core.LocalObjectReference{Name: DBSecret},
										Key:                  secretUsername,
									},
								}},
								{Name: "database_password", ValueFrom: &core.EnvVarSource{
									SecretKeyRef: &core.SecretKeySelector{
										LocalObjectReference: core.LocalObjectReference{Name: DBSecret},
										Key:                  secretPassword,
									},
								}},
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
