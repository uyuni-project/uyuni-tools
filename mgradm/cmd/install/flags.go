// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"github.com/rs/zerolog/log"
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
	if err := cmd.Flags().MarkDeprecated("tftp", "Use --tftpd-disable instead"); err != nil {
		log.Error().Err(err).Msg(L("failed to mark tftp deprecated"))
	}

	cmd_utils.AddServerFlags(cmd)

	cmd.Flags().Bool("debug-java", false, L("Enable tomcat and taskomatic remote debugging"))

	cmd_utils.AddCocoFlag(cmd)

	cmd_utils.AddHubXmlrpcFlags(cmd)

	cmd_utils.AddSalineFlag(cmd)

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "tftpd-container", Title: L("TFTPD Flags")})
	utils.AddTFTPDFlags(cmd, true, "tftpd-container")

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
