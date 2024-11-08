// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type kubernetesInstallFlags struct {
	shared.InstallFlags `mapstructure:",squash"`
	Helm                cmd_utils.HelmFlags
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[kubernetesInstallFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes [fqdn]",
		Short: L("Install a new server on a kubernetes cluster"),
		Long: L(`Install a new server on a kubernetes cluster

The install command assumes the following:
  * kubectl and helm are installed locally
  * a working kubectl configuration should be set to connect to the cluster to deploy to

The helm values file will be overridden with the values from the command parameters or configuration.

NOTE: installing on a remote cluster is not supported yet!
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetesInstallFlags
			flagsUpdater := func(v *viper.Viper) {
				flags.InstallFlags.Coco.IsChanged = v.IsSet("coco.replicas")
				flags.InstallFlags.HubXmlrpc.IsChanged = v.IsSet("hubxmlrpc.replicas")
				flags.InstallFlags.Saline.IsChanged = v.IsSet("saline.replicas") || v.IsSet("saline.port")
			}
			return utils.CommandHelper(globalFlags, cmd, args, &flags, flagsUpdater, run)
		},
	}

	shared.AddInstallFlags(cmd)
	cmd_utils.AddHelmInstallFlag(cmd)
	return cmd
}

// NewCommand for kubernetes installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, installForKubernetes)
}
