// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/migrate/shared"
	cmd_utils "github.com/uyuni-project/uyuni-tools/uyuniadm/shared/utils"
)

type podmanMigrateFlags struct {
	shared.MigrateFlags `mapstructure:",squash"`
	Podman              cmd_utils.PodmanFlags
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	migrateCmd := &cobra.Command{
		Use:   "podman [source server FQDN]",
		Short: "migrate a remote server to containers running on podman",
		Long: `Migrate a remote server to containers running on podman

This migration command assumes a few things:
  * the SSH configuration for the source server is complete, including user and
    all needed options to connect to the machine,
  * an SSH agent is started and the key to use to connect to the server is added to it,
  * podman is installed locally

NOTE: for now installing on a remote podman is not supported yet!
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viper := utils.ReadConfig(globalFlags.ConfigPath, "admconfig", cmd)
			var flags podmanMigrateFlags
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msg("Failed to Unmarshal configuration")
			}

			migrateToPodman(globalFlags, &flags, cmd, args)
		},
	}

	shared.AddMigrateFlags(migrateCmd)
	cmd_utils.AddPodmanInstallFlag(migrateCmd)

	return migrateCmd
}
