// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package rename

import (
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type renameFlags struct {
	Backend string
	SSL     adm_utils.InstallSSLFlags
}

// NewCommand creates a CLI command to rename the server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename [New FQDN]",
		Short: L("Change the host name of the server"),
		Long: L(`Set the FQDN of the server to a new value.
If no Fully Qualified Domain Name is passed, the one from the running machine will be used.

Changing the name of the server may involve updating SSL certificates to match the new name,
but also altering various configurations inside the containers.
The uyuni-server container will be stopped during the rename and
a refresh of the pillars of each registered system will be triggered.
`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags renameFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, rename)
		},
	}

	utils.AddBackendFlag(cmd)
	ssl.AddSSLGenerationFlags(cmd)
	ssl.AddSSLThirdPartyFlags(cmd)
	ssl.AddSSLDBThirdPartyFlags(cmd)

	cmd.Flags().String("ssl-password", "", L("Password for the CA key to generate"))
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-password", ssl.GeneratedFlagsGroup)
	return cmd
}

func rename(globalFlags *types.GlobalFlags, flags *renameFlags, cmd *cobra.Command, args []string) error {
	fn, err := shared.ChoosePodmanOrKubernetes(cmd.Flags(), renameForPodman, renameForKubernetes)
	if err != nil {
		return err
	}

	return fn(globalFlags, flags, cmd, args)
}
