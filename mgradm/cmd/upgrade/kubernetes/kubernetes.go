// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type kubernetesUpgradeFlags struct {
	shared.UpgradeFlags `mapstructure:",squash"`
	Helm                cmd_utils.HelmFlags
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[kubernetesUpgradeFlags]) *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "kubernetes",
		Short: L("Upgrade a local server on kubernetes"),
		Long:  L("Upgrade a local server on kubernetes"),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetesUpgradeFlags
			flagsUpdater := func(v *viper.Viper) {
				flags.UpgradeFlags.Coco.IsChanged = v.IsSet("coco.replicas")
				flags.UpgradeFlags.HubXmlrpc.IsChanged = v.IsSet("hubxmlrpc.replicas")
			}
			return utils.CommandHelper(globalFlags, cmd, args, &flags, flagsUpdater, run)
		},
	}

	shared.AddUpgradeFlags(upgradeCmd)
	cmd_utils.AddHelmInstallFlag(upgradeCmd)

	return upgradeCmd
}

// NewCommand to upgrade a kubernetes server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, upgradeKubernetes)
}
