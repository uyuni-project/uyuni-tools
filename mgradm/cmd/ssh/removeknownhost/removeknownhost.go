// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package removeknownhost

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const knownHostsPath = "/var/lib/salt/.ssh/known_hosts"

type removeKnownHostFlags struct {
	Backend string
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[removeKnownHostFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-known-host [hostname]...",
		Short: L("Remove SSH known host entries"),
		Long:  L("Remove entries from the SSH known_hosts file to avoid host key verification errors"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeKnownHostFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	utils.AddBackendFlag(cmd)
	return cmd
}

// NewCommand removes SSH known host entries from the server container.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, removeKnownHost)
}

func removeKnownHost(_ *types.GlobalFlags, flags *removeKnownHostFlags, _ *cobra.Command, args []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)

	if !cnx.TestExistenceInPod(knownHostsPath) {
		log.Info().Msg(L("No known_hosts file found, nothing to remove"))
		return nil
	}

	for _, hostname := range args {
		log.Info().Msgf(L("Removing SSH known host entry for %s"), hostname)
		if err := adm_utils.ExecCommand(
			zerolog.InfoLevel, cnx, "ssh-keygen", "-R", hostname, "-f", knownHostsPath,
		); err != nil {
			return utils.Errorf(err, L("failed to remove known host entry for %s"), hostname)
		}
	}

	log.Info().Msg(L("SSH known host entries removed successfully"))
	return nil
}
