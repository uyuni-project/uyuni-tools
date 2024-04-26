// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package support

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/support/ptf"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand to export supportconfig.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	supportCmd := &cobra.Command{
		Use:   "support",
		Short: L("Commands for support operations"),
		Long:  L("Commands for support operations"),
	}

	if ptfCommand := ptf.NewCommand(globalFlags); ptfCommand != nil {
		supportCmd.AddCommand(ptfCommand)
		return supportCmd
	}
	return nil
}
