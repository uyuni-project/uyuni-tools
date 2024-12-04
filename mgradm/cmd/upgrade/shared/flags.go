// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
)

// AddUpgradeFlags add upgrade flags to a command.
func AddUpgradeFlags(cmd *cobra.Command) {
	utils.AddImageFlag(cmd)
	utils.AddSCCFlag(cmd)
	utils.AddDBUpgradeImageFlag(cmd)

	utils.AddUpgradeCocoFlag(cmd)
	utils.AddUpgradeHubXmlrpcFlags(cmd)
	utils.AddUpgradeSalineFlag(cmd)
}

// AddUpgradeListFlags add upgrade list flags to a command.
func AddUpgradeListFlags(cmd *cobra.Command) {
	utils.AddImageFlag(cmd)
	utils.AddSCCFlag(cmd)
}
