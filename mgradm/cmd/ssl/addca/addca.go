// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package addca

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addCAFlags struct {
	SSL adm_utils.InstallSSLFlags
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[addCAFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use: "addca [fqdn]",
		// Alternative spelling; always prefer the original spelling.
		Aliases: []string{"add-ca"},
		Short:   L("Add a new root CA (phase 1 of an SSL CA rotation)"),
		Long:    L(`Add a new root CA to the trusted CA bundle without switching the server certificate.`),
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addCAFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getFlagsUpdater(&flags), run)
		},
	}

	ssl.AddSSLGenerationFlags(cmd)
	ssl.AddSSLCARootFlags(cmd)
	cmd.Flags().String("ssl-password", "", L("Password for the CA key to generate"))
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-password", ssl.GeneratedFlagsGroup)
	return cmd
}

// getFlagsUpdater defaults the database root CA from the server one when it is not provided.
func getFlagsUpdater(flags *addCAFlags) utils.FlagsUpdaterFunc {
	return func(_ *viper.Viper) {
		if flags.SSL.Ca.IsThirdParty() && !flags.SSL.DB.CA.IsThirdParty() {
			flags.SSL.DB.CA.Root = flags.SSL.Ca.Root
		}
	}
}

// NewCommand creates the command to add a new root CA to the trust (phase 1 of a rotation).
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, addCAForPodman)
}
