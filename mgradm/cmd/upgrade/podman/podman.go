// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanUpgradeFlags struct {
	cmd_utils.ServerFlags `mapstructure:",squash"`
	Podman                podman.PodmanFlags
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[podmanUpgradeFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "podman",
		Short: L("Upgrade a local server on podman"),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanUpgradeFlags
			flagsUpdater := func(v *viper.Viper) {
				flags.ServerFlags.Coco.IsChanged = v.IsSet("coco.replicas")
				flags.ServerFlags.HubXmlrpc.IsChanged = v.IsSet("hubxmlrpc.replicas")
				flags.ServerFlags.Saline.IsChanged = v.IsSet("saline.replicas") || v.IsSet("saline.port")
			}
			return utils.CommandHelper(globalFlags, cmd, args, &flags, flagsUpdater, run)
		},
	}
	shared.AddUpgradeFlags(cmd)
	podman.AddPodmanArgFlag(cmd)
	return cmd
}

func newListCmd(globalFlags *types.GlobalFlags, run func(*podmanUpgradeFlags) error) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: L("List available tags for an image"),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			viper, _ := utils.ReadConfig(cmd, utils.GlobalConfigFilename, globalFlags.ConfigPath)

			var flags podmanUpgradeFlags
			if err := viper.Unmarshal(&flags); err != nil {
				return utils.Errorf(err, L("failed to unmarshall configuration"))
			}
			if err := run(&flags); err != nil {
				return err
			}
			return nil
		},
	}
	shared.AddUpgradeListFlags(listCmd)
	return listCmd
}

func listTags(flags *podmanUpgradeFlags) error {
	hostData, err := podman.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman.PodmanLogin(hostData, flags.Image.Registry)
	if err != nil {
		return err
	}
	defer cleaner()

	return podman.ShowAvailableTag(flags.Image.Registry.Host, flags.Image, authFile)
}

// NewCommand to upgrade a podman server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := newCmd(globalFlags, upgradePodman)

	cmd.AddCommand(newListCmd(globalFlags, listTags))
	return cmd
}
