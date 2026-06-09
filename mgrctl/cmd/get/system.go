// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func newSystemCommand(globalFlags *types.GlobalFlags, parentFlags *getOptions) *cobra.Command {
	return &cobra.Command{
		Use:     "system [name/id]",
		Aliases: []string{"systems"},
		Short:   L("List or get details for registered systems"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, parentFlags, nil,
				func(_ *types.GlobalFlags, flags *getOptions, _ *cobra.Command, args []string) error {
					name := ""
					if len(args) > 0 {
						name = args[0]
					}
					return runGet(flags, systemResource{}, name)
				})
		},
	}
}
