// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
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
