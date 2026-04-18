// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package knownhost

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type apiFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func NewKnownHostCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags apiFlags

	sshKnownhostCmd := &cobra.Command{
		Use:   "knownhost",
		Short: L("SSH known_hosts file management"),
	}

	sshRemoveKnownHostCmd := &cobra.Command{
		Use:   "remove hostname [port]",
		Short: L("Remove a SSH known host"),
		Long: L(`Removes a host from the list of Salt's SSH known_hosts.
If no port is specified, it will default to 22.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, runRemoveKnownHost)
		},
		Args: cobra.RangeArgs(1, 2),
	}

	sshKnownhostCmd.AddCommand(sshRemoveKnownHostCmd)
	return sshKnownhostCmd
}
