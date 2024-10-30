// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[kubernetes.KubernetesServerFlags]) *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "kubernetes",
		Short: L("Upgrade a local server on kubernetes"),
		Long:  L("Upgrade a local server on kubernetes"),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetes.KubernetesServerFlags
			flagsUpdater := func(v *viper.Viper) {
				flags.ServerFlags.Coco.IsChanged = v.IsSet("coco.replicas")
				flags.ServerFlags.HubXmlrpc.IsChanged = v.IsSet("hubxmlrpc.replicas")
				flags.ServerFlags.Saline.IsChanged = v.IsSet("saline.replicas") || v.IsSet("saline.port")
				flags.ServerFlags.Pgsql.IsChanged = v.IsSet("pgsql.replicas")
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
