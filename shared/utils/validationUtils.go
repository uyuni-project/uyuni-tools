// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/spf13/cobra"
)

// MarkMandatoryFlags ensures that the specified flags are marked as required for the given command.
func MarkMandatoryFlags(cmd *cobra.Command, fields []string) {
	for _, field := range fields {
		if err := cmd.MarkFlagRequired(field); err != nil {
			return
		}
	}
}
