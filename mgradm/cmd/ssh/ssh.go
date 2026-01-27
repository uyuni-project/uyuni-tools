// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/ssh/removeknownhost"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand to manage SSH configuration.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	sshCmd := &cobra.Command{
		Use:     "ssh",
		GroupID: "tool",
		Short:   L("Manage SSH configuration"),
	}

	sshCmd.AddCommand(removeknownhost.NewCommand(globalFlags))

	return sshCmd
}
