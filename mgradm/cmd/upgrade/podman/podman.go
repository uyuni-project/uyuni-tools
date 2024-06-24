// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanUpgradeFlags struct {
	shared.UpgradeFlags `mapstructure:",squash"`
	Podman              podman.PodmanFlags
	MirrorPath          string
}

// NewCommand to upgrade a podman server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "podman",
		Short: L("Upgrade a local server on podman"),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanUpgradeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, upgradePodman)
		},
	}
	shared.AddUpgradeFlags(upgradeCmd)
	podman.AddPodmanArgFlag(upgradeCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: L("List available tags for an image"),
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			viper, _ := utils.ReadConfig(cmd, utils.GlobalConfigFilename, globalFlags.ConfigPath)

			var flags podmanUpgradeFlags
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msg(L("failed to unmarshall configuration"))
			}
			tags, _ := podman.ShowAvailableTag(globalFlags.Registry, flags.Image)
			log.Info().Msgf(L("Available Tags for image: %s"), flags.Image.Name)
			for _, value := range tags {
				log.Info().Msgf(value)
			}
		},
	}
	shared.AddUpgradeListFlags(listCmd)
	upgradeCmd.AddCommand(listCmd)

	return upgradeCmd
}
