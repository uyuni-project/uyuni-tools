// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package upgrade

import (
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// AddUpgradeFlags add upgrade flags to a command.
func AddUpgradeFlags(cmd *cobra.Command) {
	adm_utils.AddServerFlags(cmd)

	adm_utils.AddDBUpgradeImageFlag(cmd)
	adm_utils.AddUpgradeCocoFlag(cmd)
	adm_utils.AddUpgradeHubXmlrpcFlags(cmd)
	adm_utils.AddUpgradeSalineFlag(cmd)
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "tftpd-container", Title: L("TFTPD Flags")})
	utils.AddTFTPDFlags(cmd, true, "tftpd-container")
}

// AddUpgradeListFlags add upgrade list flags to a command.
func AddUpgradeListFlags(cmd *cobra.Command) {
	adm_utils.AddImageFlag(cmd)
	adm_utils.AddSCCFlag(cmd)
}
