// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package org

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand  command for APIs.
func NewCommand(globalFlags *types.GlobalFlags) (*cobra.Command, error) {
	orgCmd := &cobra.Command{
		Use:   "org",
		Short: L("Organization-related commands"),
	}

	if err := api.AddAPIFlags(orgCmd, false); err != nil {
		return orgCmd, err
	}

	orgCmd.AddCommand(createFirstCommand(globalFlags))

	return orgCmd, nil
}
