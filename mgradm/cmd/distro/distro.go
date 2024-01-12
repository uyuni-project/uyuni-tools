// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Backend           string
	ChannelLabel      string
	ProductMap        map[string]map[string]types.Distribution
	ConnectionDetails api.ConnectionDetails `mapstructure:"api"`
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags flagpole

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
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, distroCp)
		},
	}

	api.AddAPIFlags(distroCmd, true)
	distroCmd.AddCommand(cpCmd)
	return distroCmd
}
