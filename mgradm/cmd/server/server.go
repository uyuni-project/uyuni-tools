// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/server/rename"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand creates a sub command for all server-related actions.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server",
		GroupID: "tool",
		Short:   L("Server management utilities"),
		Args:    cobra.ExactArgs(1),
	}

	cmd.AddCommand(rename.NewCommand(globalFlags))
	return cmd
}
