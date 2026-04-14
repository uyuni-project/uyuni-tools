// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/ssh/knownhost"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	sshCmd := &cobra.Command{
		Use:   "ssh",
		Short: L("SSH management commands"),
	}

	sshCmd.AddCommand(knownhost.NewKnownHostCommand(globalFlags))
	api.AddAPIFlags(sshCmd)

	return sshCmd
}
