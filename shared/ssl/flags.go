// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package ssl

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func AddSSLGenerationFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice("ssl-cname", []string{}, L("SSL certificate cnames separated by commas"))
	cmd.Flags().String("ssl-country", "DE", L("SSL certificate country"))
	cmd.Flags().String("ssl-state", "Bayern", L("SSL certificate state"))
	cmd.Flags().String("ssl-city", "Nuernberg", L("SSL certificate city"))
	cmd.Flags().String("ssl-org", "SUSE", L("SSL certificate organization"))
	cmd.Flags().String("ssl-ou", "SUSE", L("SSL certificate organization unit"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "ssl", Title: L("SSL Certificate Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-cname", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-country", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-state", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-city", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-org", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ou", "ssl")
}
