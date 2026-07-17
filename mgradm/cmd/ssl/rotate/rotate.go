// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package rotate

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type rotateFlags struct {
	SSL       adm_utils.InstallSSLFlags
	Force     bool
	Emergency bool
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[rotateFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate [fqdn]",
		Short: L("Switch to a server certificate signed by the new CA, then drop the old CA"),
		Long: L(`Second phase of an SSL CA rotation, to run after 'mgradm ssl addca' has added the new CA and the
clients trust it: issue a new server (and database) certificate signed by the new CA, switch the
services to it, and drop the old CA from the trusted bundle so only the new CA is left.`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags rotateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getFlagsUpdater(&flags), run)
		},
	}

	ssl.AddSSLGenerationFlags(cmd)
	ssl.AddSSLThirdPartyFlags(cmd)
	ssl.AddSSLDBThirdPartyFlags(cmd)
	cmd.Flags().String("ssl-password", "", L("Password for the CA key"))
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-password", ssl.GeneratedFlagsGroup)
	cmd.Flags().Bool("force", false, L("Rotate even if some clients do not trust the new CA yet"))
	cmd.Flags().Bool("check-only", false, L("Only report client readiness, without rotating"))
	cmd.Flags().Bool("emergency", false,
		L("Generate a new CA and drop the old one immediately without an overlap period (e.g. for compromised CAs)"))
	return cmd
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, rotateForPodman)
}

// getFlagsUpdater defaults the database SSL flags from the server ones when they are not provided.
func getFlagsUpdater(flags *rotateFlags) utils.FlagsUpdaterFunc {
	return func(_ *viper.Viper) {
		if flags.SSL.Ca.IsThirdParty() && !flags.SSL.DB.CA.IsThirdParty() {
			flags.SSL.DB.CA.Root = flags.SSL.Ca.Root
			flags.SSL.DB.CA.Intermediate = flags.SSL.Ca.Intermediate
		}
		if flags.SSL.Server.IsDefined() && !flags.SSL.DB.IsDefined() {
			flags.SSL.DB.Cert = flags.SSL.Server.Cert
			flags.SSL.DB.Key = flags.SSL.Server.Key
		}
	}
}
