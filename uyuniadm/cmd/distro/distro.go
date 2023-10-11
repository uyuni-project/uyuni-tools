// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Backend      string
	ChannelLabel string
	ProductMap   map[string]map[string]types.Distribution
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}
	apiFlags := &api.ConnectionDetails{}

	distroCmd := &cobra.Command{
		Use:     "distribution",
		Short:   "Distribution management",
		Long:    "Tools and utilities for distribution management",
		Aliases: []string{"distro"},
	}

	cpCmd := &cobra.Command{
		Use:   "copy path-to-source distribution-name [channel-label]",
		Short: "Copy distribution files from iso to the container",
		Long: `Takes a path to iso file or directory with mounted iso and copies it into the container.

Distribution name specifies the destination directory under /srv/www/distributions.

Optional channel label specify which parent channel to associate with the distribution. Only when API details are provided and auto registration is done.`,
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"cp"},
		Run: func(cmd *cobra.Command, args []string) {
			viper := utils.ReadConfig(globalFlags.ConfigPath, "admconfig", cmd)
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshall configuration")
			}
			var channelLabel string
			if len(args) == 3 {
				channelLabel = args[2]
			} else {
				channelLabel = ""
			}
			distCp(globalFlags, flags, apiFlags, cmd, args[1], args[0], channelLabel)
		},
	}

	api.AddAPIFlags(distroCmd, apiFlags, true)
	distroCmd.AddCommand(cpCmd)
	return distroCmd
}
