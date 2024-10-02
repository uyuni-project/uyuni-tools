// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
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
			flags.ServerFlags.Coco.IsChanged = cmd.Flags().Changed("coco-replicas")
			flags.ServerFlags.HubXmlrpc.IsChanged = cmd.Flags().Changed("hubxmlrpc-replicas")
			return utils.CommandHelper(globalFlags, cmd, args, &flags, run)
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
