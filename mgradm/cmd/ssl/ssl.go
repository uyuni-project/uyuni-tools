// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package ssl

import (
	"github.com/spf13/cobra"
	ssladdca "github.com/uyuni-project/uyuni-tools/mgradm/cmd/ssl/addca"
	sslrotate "github.com/uyuni-project/uyuni-tools/mgradm/cmd/ssl/rotate"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand creates the command to rotate the server SSL CA and certificates.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	sslCmd := &cobra.Command{
		Use:     "ssl",
		GroupID: "tool",
		Short:   L("Rotate the server SSL CA and certificates"),
		Long:    L("Rotate the server and database SSL CA and certificates."),
		Args:    cobra.ExactArgs(1),
	}

	sslCmd.AddCommand(ssladdca.NewCommand(globalFlags))
	sslCmd.AddCommand(sslrotate.NewCommand(globalFlags))

	return sslCmd
}
