// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/spf13/cobra"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// AddInstallFlags add flags to installa command.
func AddInstallFlags(cmd *cobra.Command) {
	cmd.Flags().String("tz", "", L("Time zone to set on the server. Defaults to the host timezone"))
	cmd.Flags().String("email", "admin@example.com", L("Administrator e-mail"))
	cmd.Flags().String("emailfrom", "notifications@example.com", L("E-Mail sending the notifications"))
	cmd.Flags().String("issParent", "", L("InterServerSync v1 parent FQDN"))

	cmd_utils.AddDBFlags(cmd)
	cmd_utils.AddReportDBFlags(cmd)

	cmd.Flags().Bool("tftp", true, L("Enable TFTP"))

	// For generated CA and certificate
	ssl.AddSSLGenerationFlags(cmd)
	cmd.Flags().String("ssl-password", "", L("Password for the CA key to generate"))
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-password", "ssl")

	// For SSL 3rd party certificates
	cmd.Flags().StringSlice("ssl-ca-intermediate", []string{}, L("Intermediate CA certificate path"))
	cmd.Flags().String("ssl-ca-root", "", L("Root CA certificate path"))
	cmd.Flags().String("ssl-server-cert", "", L("Server certificate path"))
	cmd.Flags().String("ssl-server-key", "", L("Server key path"))

	cmd.Flags().StringSlice("ssl-db-ca-intermediate", []string{},
		L("Intermediate CA certificate path for the database if different from the server one"))
	cmd.Flags().String("ssl-db-ca-root", "",
		L("Root CA certificate path for the database if different from the server one"))
	cmd.Flags().String("ssl-db-cert", "", L("Database certificate path"))
	cmd.Flags().String("ssl-db-key", "", L("Database key path"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "ssl3rd", Title: L("3rd Party SSL Certificate Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ca-intermediate", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ca-root", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-server-cert", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-server-key", "ssl3rd")

	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-ca-intermediate", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-ca-root", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-cert", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-key", "ssl3rd")

	cmd_utils.AddSCCFlag(cmd)

	cmd.Flags().Bool("debug-java", false, L("Enable tomcat and taskomatic remote debugging"))
	cmd_utils.AddImageFlag(cmd)

	cmd_utils.AddCocoFlag(cmd)

	cmd_utils.AddHubXmlrpcFlags(cmd)

	cmd_utils.AddSalineFlag(cmd)

	cmd_utils.AddPgsqlFlags(cmd)

	cmd.Flags().String("admin-login", "admin", L("Administrator user name"))
	cmd.Flags().String("admin-password", "", L("Administrator password"))
	cmd.Flags().String("admin-firstName", "Administrator", L("First name of the administrator"))
	cmd.Flags().String("admin-lastName", "McAdmin", L("Last name of the administrator"))
	cmd.Flags().String("organization", "Organization", L("First organization name"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "first-user", Title: L("First User Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "admin-login", "first-user")
	_ = utils.AddFlagToHelpGroupID(cmd, "admin-password", "first-user")
	_ = utils.AddFlagToHelpGroupID(cmd, "admin-firstName", "first-user")
	_ = utils.AddFlagToHelpGroupID(cmd, "admin-lastName", "first-user")
	_ = utils.AddFlagToHelpGroupID(cmd, "organization", "first-user")
}
