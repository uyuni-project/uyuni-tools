// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package support

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/support/config"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand to export supportconfig.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	supportCmd := &cobra.Command{
		Use:   "support",
		Short: "commands for support operations",
		Long:  "Commands for support operations",
	}
	supportCmd.AddCommand(config.NewCommand(globalFlags))

	return supportCmd
}
