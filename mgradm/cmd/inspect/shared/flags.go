// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/spf13/cobra"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// InspectFlags are the flags used by inspect commands.
type InspectFlags struct {
	Image types.ImageFlags `mapstructure:",squash"`
	SCC   types.SCCCredentials
}

// AddInspectFlags add flags to inspect command.
func AddInspectFlags(cmd *cobra.Command) {
	cmd_utils.AddSCCFlag(cmd)
	cmd_utils.AddImageFlag(cmd)
}
