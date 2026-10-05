// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package gpg

import (
	"github.com/spf13/cobra"
	gpgadd "github.com/uyuni-project/uyuni-tools/mgradm/cmd/gpg/add"
	gpglist "github.com/uyuni-project/uyuni-tools/mgradm/cmd/gpg/list"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand import gpg keys from 3rd party repository.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	gpgKeyCmd := &cobra.Command{
		Use:        "gpg",
		GroupID:    "tool",
		Short:      L("Manage GPG keys for 3rd party repositories"),
		Args:       cobra.ExactArgs(1),
		Deprecated: "please use `mgrctl gpg` instead",
	}

	gpgKeyCmd.AddCommand(gpgadd.NewCommand(globalFlags))
	gpgKeyCmd.AddCommand(gpglist.NewCommand(globalFlags))

	return gpgKeyCmd
}
