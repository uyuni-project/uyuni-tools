// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build !ptf

package ptf

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand is the command for creates supportptf.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return nil
}
