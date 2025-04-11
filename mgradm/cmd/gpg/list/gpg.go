// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package gpglist

import (
	"strings"

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

const customKeyringPath = "/var/spacewalk/gpg/customer-build-keys.gpg"
const systemKeyringPath = "/usr/lib/susemanager/susemanager-build-keys.gpg"

type gpgListFlags struct {
	Backend string
	System  bool
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[gpgListFlags]) *cobra.Command {
	gpgListKeyCmd := &cobra.Command{
		Use:   "list",
		Short: L("List GPG keys"),
		Long:  L("List GPG keys from custom keyring (default) or system keyring"),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags gpgListFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	gpgListKeyCmd.Flags().BoolP("system", "s", false, L("List keys from system keyring"))
	utils.AddBackendFlag(gpgListKeyCmd)
	return gpgListKeyCmd
}

// NewCommand lists imported gpg keys.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, gpgListKeys)
}

func gpgListKeys(_ *types.GlobalFlags, flags *gpgListFlags, _ *cobra.Command, _ []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)

	gpgListCmd := []string{"gpg", "--no-default-keyring", "--keyring"}

	if flags.System {
		gpgListCmd = append(gpgListCmd, systemKeyringPath, "--list-keys")
	} else {
		gpgListCmd = append(gpgListCmd, customKeyringPath, "--list-keys")
	}

	log.Info().Msgf(L("Running %s"), strings.Join(gpgListCmd, " "))
	if err := adm_utils.ExecCommand(
		zerolog.InfoLevel, cnx, gpgListCmd...,
	); err != nil {
		return utils.Errorf(err, L("failed to list keys in selected keyring"))
	}

	return nil
}
