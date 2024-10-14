// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/kubernetes"
	pxy_utils "github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type kubernetesPTFFlags struct {
	UpgradeFlags kubernetes.KubernetesProxyUpgradeFlags `mapstructure:",squash"`
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[kubernetesPTFFlags]) *cobra.Command {
	kubernetesCmd := &cobra.Command{
		Use:   "kubernetes",
		Short: L("Install a PTF or Test package on a kubernetes cluster"),
		Long: L(`Install a PTR of Test package on a kubernetes cluster

The support ptf command assumes the following:
  * kubectl and helm are installed locally
  * a working kubectl configuration should be set to connect to the cluster to deploy to

The helm values file will be overridden with the values from the command parameters or configuration.

NOTE: installing on a remote cluster is not supported yet!
`),

		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetesPTFFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	pxy_utils.AddImageFlags(kubernetesCmd)

	kubernetes.AddHelmFlags(kubernetesCmd)

	return kubernetesCmd
}

// NewCommand for kubernetes installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, ptfForKubernetes)
}
