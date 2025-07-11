// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func upgradeKubernetes(_ *types.GlobalFlags,
	flags *kubernetes.KubernetesProxyUpgradeFlags, cmd *cobra.Command, args []string,
) error {
	flags.ProxyImageFlags.CheckParameters()
	return kubernetes.Upgrade(flags, cmd, args)
}
