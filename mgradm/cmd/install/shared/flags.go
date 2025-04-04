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
	cmd.Flags().String("email", "admin@example.com", L("Administrator e-mail"))
	cmd.Flags().String("emailfrom", "notifications@example.com", L("E-Mail sending the notifications"))
	cmd.Flags().String("issParent", "", L("InterServerSync v1 parent FQDN"))
	cmd.Flags().Bool("tftp", true, L("Enable TFTP"))

	cmd.Flags().Bool("debug-java", false, L("Enable tomcat and taskomatic remote debugging"))

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

	cmd_utils.AddServerFlags(cmd)

	cmd_utils.AddInstallSSLFlag(cmd)
}
