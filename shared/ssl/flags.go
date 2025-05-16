// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package ssl

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const (
	GeneratedFlagsGroup  = "ssl"
	ThirdPartyFlagsGroup = "ssl3rd"
)

// AddSSLGenerationFlags adds the command flags to generate SSL certificates.
func AddSSLGenerationFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice("ssl-cname", []string{}, L("SSL certificate cnames separated by commas"))
	cmd.Flags().String("ssl-country", "DE", L("SSL certificate country"))
	cmd.Flags().String("ssl-state", "Bayern", L("SSL certificate state"))
	cmd.Flags().String("ssl-city", "Nuernberg", L("SSL certificate city"))
	cmd.Flags().String("ssl-org", "SUSE", L("SSL certificate organization"))
	cmd.Flags().String("ssl-ou", "SUSE", L("SSL certificate organization unit"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: GeneratedFlagsGroup, Title: L("SSL Certificate Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-cname", GeneratedFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-country", GeneratedFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-state", GeneratedFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-city", GeneratedFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-org", GeneratedFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ou", GeneratedFlagsGroup)
}

// AddSSLThirdPartyFlags adds the command flags to pass Apache third party SSL certificates.
func AddSSLThirdPartyFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice("ssl-ca-intermediate", []string{}, L("Intermediate CA certificate path"))
	cmd.Flags().String("ssl-ca-root", "", L("Root CA certificate path"))
	cmd.Flags().String("ssl-server-cert", "", L("Server certificate path"))
	cmd.Flags().String("ssl-server-key", "", L("Server key path"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: ThirdPartyFlagsGroup, Title: L("3rd Party SSL Certificate Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ca-intermediate", ThirdPartyFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ca-root", ThirdPartyFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-server-cert", ThirdPartyFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-server-key", ThirdPartyFlagsGroup)
}

// AddSSLDBThirdPartyFlags adds the command flags to pass database third party SSL certificates.
func AddSSLDBThirdPartyFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice("ssl-db-ca-intermediate", []string{},
		L("Intermediate CA certificate path for the database if different from the server one"))
	cmd.Flags().String("ssl-db-ca-root", "",
		L("Root CA certificate path for the database if different from the server one"))
	cmd.Flags().String("ssl-db-cert", "", L("Database certificate path"))
	cmd.Flags().String("ssl-db-key", "", L("Database key path"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: ThirdPartyFlagsGroup, Title: L("3rd Party SSL Certificate Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-ca-intermediate", ThirdPartyFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-ca-root", ThirdPartyFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-cert", ThirdPartyFlagsGroup)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-db-key", ThirdPartyFlagsGroup)
}
