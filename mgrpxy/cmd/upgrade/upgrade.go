// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package upgrade

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/upgrade/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/upgrade/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand install a new proxy from scratch.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: L("Upgrade a proxy"),
		Long:  L("Upgrade a proxy"),
	}
	upgradeCmd.PersistentFlags().StringVar(&globalFlags.Registry, "registry", "", L("specify a private registry"))

	upgradeCmd.AddCommand(podman.NewCommand(globalFlags))
	upgradeCmd.AddCommand(kubernetes.NewCommand(globalFlags))

	return upgradeCmd
}
