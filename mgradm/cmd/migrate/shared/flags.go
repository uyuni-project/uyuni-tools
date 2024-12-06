// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// AddMigrateFlags add migration flags to a command.
func AddMigrateFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("prepare", false, L("Prepare the mgration - copy the data without stopping the source server."))
	utils.AddMirrorFlag(cmd)
	utils.AddSCCFlag(cmd)
	utils.AddImageFlag(cmd)
	utils.AddDBUpgradeImageFlag(cmd)
	utils.AddUpgradeCocoFlag(cmd)
	utils.AddUpgradeHubXmlrpcFlags(cmd)
	utils.AddUpgradeSalineFlag(cmd)
	cmd.Flags().String("user", "root",
		L("User on the source server. Non-root user must have passwordless sudo privileges (NOPASSWD tag in /etc/sudoers)."),
	)
}
