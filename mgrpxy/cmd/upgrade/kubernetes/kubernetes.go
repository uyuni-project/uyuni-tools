// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/kubernetes"
	pxy_utils "github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func newCmd(
	globalFlags *types.GlobalFlags,
	run utils.CommandFunc[kubernetes.KubernetesProxyUpgradeFlags],
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes",
		Short: L("Upgrade a proxy on a running kubernetes cluster"),
		Long: L(`Upgrade a proxy on a running kubernetes cluster.

The upgrade kubernetes command assumes kubectl is installed locally.

NOTE: for now upgrading on a remote kubernetes cluster is not supported!
`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetes.KubernetesProxyUpgradeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	pxy_utils.AddImageFlags(cmd)

	kubernetes.AddHelmFlags(cmd)

	return cmd
}

// NewCommand install a new proxy on a running kubernetes cluster.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, upgradeKubernetes)
}
