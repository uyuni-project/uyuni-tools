// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package org

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	orgCmd := &cobra.Command{
		Use:   "org",
		Short: "Organization-related commands",
	}

	api.AddAPIFlags(orgCmd, false)

	orgCmd.AddCommand(createFirstCommand(globalFlags))

	return orgCmd
}
