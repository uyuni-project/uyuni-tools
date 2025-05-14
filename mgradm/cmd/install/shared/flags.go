// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/spf13/cobra"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// AddInstallFlags add flags to install command.
func AddInstallFlags(cmd *cobra.Command) {
	cmd.Flags().String("tz", "", L("Time zone to set on the server. Defaults to the host timezone"))
	cmd.Flags().String("email", "admin@example.com", L("Administrator e-mail"))
	cmd.Flags().String("emailfrom", "notifications@example.com", L("E-Mail sending the notifications"))
	cmd.Flags().String("issParent", "", L("InterServerSync v1 parent FQDN"))
	cmd.Flags().Bool("tftp", true, L("Enable TFTP"))

	cmd_utils.AddServerFlags(cmd)

	// For SSL 3rd party certificates
	cmd.Flags().StringSlice("ssl-ca-intermediate", []string{}, L("Intermediate CA certificate path"))
	cmd.Flags().String("ssl-ca-root", "", L("Root CA certificate path"))
	cmd.Flags().String("ssl-server-cert", "", L("Server certificate path"))
	cmd.Flags().String("ssl-server-key", "", L("Server key path"))

	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ca-intermediate", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ca-root", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-server-cert", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-server-key", "ssl3rd")

	cmd.Flags().Bool("debug-java", false, L("Enable tomcat and taskomatic remote debugging"))

	cmd_utils.AddCocoFlag(cmd)

	cmd_utils.AddHubXmlrpcFlags(cmd)

	cmd_utils.AddSalineFlag(cmd)

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
