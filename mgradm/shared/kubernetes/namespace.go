// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CreateNamespace creates a kubernetes namespace.
func CreateNamespace(namespace string) error {
	ns := core.Namespace{
		TypeMeta: meta.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
		ObjectMeta: meta.ObjectMeta{
			Name: namespace,
		},
	}
	return kubernetes.Apply([]runtime.Object{&ns}, L("failed to create the namespace"))
}
