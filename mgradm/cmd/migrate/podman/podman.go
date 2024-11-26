// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	podman_utils "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanMigrateFlags struct {
	shared.MigrateFlags `mapstructure:",squash"`
	SCC                 types.SCCCredentials
	Podman              podman_utils.PodmanFlags
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[podmanMigrateFlags]) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "podman [source server FQDN]",
		Short: L("Migrate a remote server to containers running on podman"),
		Long: L(`Migrate a remote server to containers running on podman

This migration command assumes a few things:
  * the SSH configuration for the source server is complete, including user and
    all needed options to connect to the machine,
  * an SSH agent is started and the key to use to connect to the server is added to it,
  * podman is installed locally

NOTE: migrating to a remote podman is not supported yet!
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanMigrateFlags
			flagsUpdater := func(v *viper.Viper) {
				flags.MigrateFlags.Coco.IsChanged = v.IsSet("coco.replicas")
				flags.MigrateFlags.HubXmlrpc.IsChanged = v.IsSet("hubxmlrpc.replicas")
				flags.MigrateFlags.Saline.IsChanged = v.IsSet("saline.replicas") || v.IsSet("saline.port")
			}
			return utils.CommandHelper(globalFlags, cmd, args, &flags, flagsUpdater, run)
		},
	}

	shared.AddMigrateFlags(migrateCmd)
	podman_utils.AddPodmanArgFlag(migrateCmd)

	return migrateCmd
}

// NewCommand for podman migration.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, migrateToPodman)
}
