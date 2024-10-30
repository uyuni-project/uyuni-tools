// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	// DBDeployName is the name of the database Deployment object.
	DBDeployName = "db"
	// DBAdminSecret is the name of the database administrator credentials secret.
	// This secret is only needed for a DB prepared by mgradm.
	DBAdminSecret = "db-admin-credentials"
	// DBSecret is the name of the database credentials secret.
	DBSecret = "db-credentials"
	// ReportdbSecret is the name of the report database credentials secret.
	ReportdbSecret = "reportdb-credentials"
	SCCSecret      = "scc-credentials"
	secretUsername = "username"
	secretPassword = "password"
)

// CreateBasicAuthSecret creates a secret of type basic-auth.
func CreateBasicAuthSecret(namespace string, name string, user string, password string) error {
	// Check if the secret is already existing
	out, err := runCmdOutput(zerolog.DebugLevel, "kubectl", "get", "-n", namespace, "secret", name, "-o", "name")
	if err == nil && strings.TrimSpace(string(out)) != "" {
		return nil
	}

	// Create the secret
	secret := core.Secret{
		TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "Secret"},
		ObjectMeta: meta.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels:    kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.ServerComponent),
		},
		// It seems serializing this object automatically transforms the secrets to base64.
		Data: map[string][]byte{
			secretUsername: []byte(user),
			secretPassword: []byte(password),
		},
		Type: core.SecretTypeBasicAuth,
	}

	return kubernetes.Apply([]runtime.Object{&secret}, L("failed to create the secret"))
}

// CreateDBDeployment creates a new deployment of the database.
func CreateDBDeployment(
	namespace string,
	image string,
	pullPolicy string,
	pullSecret string,
	timezone string,
) error {
	deploy := getDBDeployment(namespace, image, kubernetes.GetPullPolicy(pullPolicy), pullSecret, timezone)
	return kubernetes.Apply([]runtime.Object{deploy}, L("failed to create the database deployment"))
}

func getDBDeployment(
	namespace string,
	image string,
	pullPolicy core.PullPolicy,
	pullSecret string,
	timezone string,
) *apps.Deployment {
	var replicas int32 = 1

	mounts := []types.VolumeMount{utils.VarPgsqlDataVolumeMount}
	volumeMounts := kubernetes.ConvertVolumeMounts(mounts)
	volumes := kubernetes.CreateVolumes(mounts)

	// Add TLS secret
	const tlsVolumeName = "tls-secret"
	var secretMode int32 = 0400
	tlsVolume := kubernetes.CreateSecretVolume(tlsVolumeName, kubernetes.DBCertSecretName)
	tlsVolume.Secret.Items = []core.KeyToPath{
		{Key: "tls.crt", Path: "tls/certs/spacewalk.crt"},
		{Key: "tls.key", Path: "tls/private/pg-spacewalk.key", Mode: &secretMode},
		{Key: "ca.crt", Path: "trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT"},
	}
	volumes = append(volumes, tlsVolume)
	volumeMounts = append(volumeMounts,
		core.VolumeMount{Name: tlsVolumeName, MountPath: "/etc/pki"},
	)

	envs := []core.EnvVar{
		{Name: "TZ", Value: timezone},
		// Add the admin credentials secret
		{
			Name: "POSTGRES_USER",
			ValueFrom: &core.EnvVarSource{
				SecretKeyRef: &core.SecretKeySelector{
					LocalObjectReference: core.LocalObjectReference{Name: DBAdminSecret}, Key: secretUsername,
				},
			},
		},
		{
			Name: "POSTGRES_PASSWORD",
			ValueFrom: &core.EnvVarSource{
				SecretKeyRef: &core.SecretKeySelector{
					LocalObjectReference: core.LocalObjectReference{Name: DBAdminSecret}, Key: secretPassword,
				},
			},
		},
		// Add the internal db user credentials secret
		{
			Name: "MANAGER_USER",
			ValueFrom: &core.EnvVarSource{
				SecretKeyRef: &core.SecretKeySelector{
					LocalObjectReference: core.LocalObjectReference{Name: DBSecret}, Key: secretUsername,
				},
			},
		},
		{
			Name: "MANAGER_PASS",
			ValueFrom: &core.EnvVarSource{
				SecretKeyRef: &core.SecretKeySelector{
					LocalObjectReference: core.LocalObjectReference{Name: DBSecret}, Key: secretPassword,
				},
			},
		},
		// Add the report db user credentials secret
		{
			Name: "REPORT_DB_USER",
			ValueFrom: &core.EnvVarSource{
				SecretKeyRef: &core.SecretKeySelector{
					LocalObjectReference: core.LocalObjectReference{Name: ReportdbSecret}, Key: secretUsername,
				},
			},
		},
		{
			Name: "REPORT_DB_PASS",
			ValueFrom: &core.EnvVarSource{
				SecretKeyRef: &core.SecretKeySelector{
					LocalObjectReference: core.LocalObjectReference{Name: ReportdbSecret}, Key: secretPassword,
				},
			},
		},
	}

	// fsGroup is required to set the owner of the mounted files, most importantly the SSL key file.
	var fsGroup int64 = 999

	deploy := apps.Deployment{
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      DBDeployName,
			Namespace: namespace,
			Labels:    kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.DBComponent),
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			// Since the DB container will never be able to scale, we need to stick to recreate strategy
			// or the new deployed pods won't be ready.
			Strategy: apps.DeploymentStrategy{Type: apps.RecreateDeploymentStrategyType},
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{kubernetes.ComponentLabel: kubernetes.DBComponent},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: kubernetes.GetLabels(kubernetes.ServerApp, kubernetes.DBComponent),
				},
				Spec: core.PodSpec{
					SecurityContext: &core.PodSecurityContext{
						FSGroup: &fsGroup,
					},
					Containers: []core.Container{
						{
							Name:            "db",
							Image:           image,
							ImagePullPolicy: pullPolicy,
							VolumeMounts:    volumeMounts,
							Env:             envs,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}

	if pullSecret != "" {
		deploy.Spec.Template.Spec.ImagePullSecrets = []core.LocalObjectReference{{Name: pullSecret}}
	}

	return &deploy
}
