// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// AddMigrateFlags add migration flags to a command.
func AddMigrateFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("prepare", false, L("Prepare the mgration - copy the data without stopping the source server."))

	cmd.Flags().String("user", "root",
		L("User on the source server. Non-root user must have passwordless sudo privileges (NOPASSWD tag in /etc/sudoers)."),
	)

	adm_utils.AddImageFlag(cmd)
	adm_utils.AddSCCFlag(cmd)

	// For generated CA and certificate
	ssl.AddSSLGenerationFlags(cmd)
	cmd.Flags().String("ssl-password", "", L("Password for the CA key to generate"))
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-password", "ssl")

	// For SSL 3rd party certificates
	cmd.Flags().StringSlice("ssl-db-ca-intermediate", []string{},
		L("Intermediate CA certificate path for the database"))
	cmd.Flags().String("ssl-db-ca-root", "",
		L("Root CA certificate path for the database"))
	cmd.Flags().String("ssl-db-cert", "", L("Database certificate path"))
	cmd.Flags().String("ssl-db-key", "", L("Database key path"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "ssl3rd", Title: L("3rd Party SSL Certificate Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-ca-intermediate", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-ca-root", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-cert", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-key", "ssl3rd")

	adm_utils.AddDBFlags(cmd)
	adm_utils.AddReportDBFlags(cmd)

	adm_utils.AddDBUpgradeImageFlag(cmd)
	adm_utils.AddUpgradeCocoFlag(cmd)
	adm_utils.AddUpgradeHubXmlrpcFlags(cmd)
	adm_utils.AddUpgradeSalineFlag(cmd)
	adm_utils.AddPgsqlFlags(cmd)
}
