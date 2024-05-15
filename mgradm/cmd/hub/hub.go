// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package hub

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/hub/register"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand command for Hub management.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	hubCmd := &cobra.Command{
		Use:     "hub",
		GroupID: "management",
		Short:   L("Hub management"),
		Long:    L("Tools and utilities for Hub management"),
		Aliases: []string{"hub"},
	}

	hubCmd.SetUsageTemplate(hubCmd.UsageTemplate())
	hubCmd.AddCommand(register.NewCommand(globalFlags))
	return hubCmd
}
