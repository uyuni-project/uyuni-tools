// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "github.com/spf13/cobra"

// AddUninstallFlags adds the common flags for uninstall commands.
func ValidateMandatoryFlags(cmd *cobra.Command, fields []string) {
	for _, field := range fields {
		if err := cmd.MarkFlagRequired(field); err != nil {
			return
		}
	}
}
