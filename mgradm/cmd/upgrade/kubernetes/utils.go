// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func upgradeKubernetes(
	_ *types.GlobalFlags,
	flags *kubernetes.KubernetesServerFlags,
	cmd *cobra.Command,
	_ []string,
) error {
	flags.Installation.CheckUpgradeParameters(cmd, "kubectl")
	return kubernetes.Reconcile(flags, "")
}
