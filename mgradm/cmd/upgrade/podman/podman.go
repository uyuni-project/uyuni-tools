// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanUpgradeFlags struct {
	shared.UpgradeFlags `mapstructure:",squash"`
	SCC                 types.SCCCredentials
	Podman              podman.PodmanFlags
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[podmanUpgradeFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "podman",
		Short: L("Upgrade a local server on podman"),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanUpgradeFlags
			flagsUpdater := func(v *viper.Viper) {
				flags.UpgradeFlags.Coco.IsChanged = v.IsSet("coco.replicas")
				flags.UpgradeFlags.HubXmlrpc.IsChanged = v.IsSet("hubxmlrpc.replicas")
			}
			return utils.CommandHelper(globalFlags, cmd, args, &flags, flagsUpdater, run)
		},
	}
	shared.AddUpgradeFlags(cmd)
	podman.AddPodmanArgFlag(cmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: L("List available tags for an image"),
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			viper, _ := utils.ReadConfig(cmd, utils.GlobalConfigFilename, globalFlags.ConfigPath)

			var flags podmanUpgradeFlags
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msg(L("failed to unmarshall configuration"))
			}
			if err := podman.ShowAvailableTag(flags.Image.Registry, flags.Image); err != nil {
				log.Fatal().Err(err)
			}
		},
	}
	shared.AddUpgradeListFlags(listCmd)
	cmd.AddCommand(listCmd)
	return cmd
}

// NewCommand to upgrade a podman server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, upgradePodman)
}
