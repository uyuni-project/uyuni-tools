// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Backend           string
	ChannelLabel      string
	ProductMap        map[string]map[string]types.Distribution
	ConnectionDetails api.ConnectionDetails `mapstructure:"api"`
}

// NewCommand command for distribution management.
func NewCommand(globalFlags *types.GlobalFlags) (*cobra.Command, error) {
	var flags flagpole

	distroCmd := &cobra.Command{
		Use:     "distribution",
		Short:   L("Distributions management"),
		Long:    L("Tools for autoinstallation distributions management"),
		Aliases: []string{"distro"},
	}

	cpCmd := &cobra.Command{
		Use:   "copy path-to-source distribution-name [channel-label]",
		Short: L("Copy distribution files from ISO image to the container"),
		Long: L(`Takes a path to an ISO file or the directory of a mounted ISO image and copies it into the container.

Distribution name specifies the destination directory under /srv/www/distributions.

Optional channel label specify which parent channel to associate with the distribution.
Only when API informations are provided and auto registration is done.`),
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"cp"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, distroCp)
		},
	}

	if err := api.AddAPIFlags(distroCmd, true); err != nil {
		return distroCmd, err
	}
	distroCmd.AddCommand(cpCmd)
	return distroCmd, nil
}
