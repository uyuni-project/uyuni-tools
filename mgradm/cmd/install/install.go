// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanInstallFlags struct {
	adm_utils.ServerFlags `mapstructure:",squash"`
	Podman                podman.PodmanFlags
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[podmanInstallFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install [fqdn]",
		Aliases: []string{"install podman"},
		GroupID: "deploy",
		Short:   L("Install a new server on podman"),
		Long: L(`Install a new server on podman

The command assumes podman is installed locally.

NOTE: installing on a remote podman is not supported yet!
`),
		Args: func(cmd *cobra.Command, args []string) error {
			// If the alias "install podman" is used, "podman" will be the first arg.
			// We remove it from the args slice so it isn't treated as the FQDN.
			if len(args) > 0 && args[0] == "podman" {
				copy(args, args[1:])
				args = args[:len(args)-1]
			}
			return cobra.MaximumNArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanInstallFlags
			flagsUpdater := func(v *viper.Viper) {
				flags.Coco.IsChanged = v.IsSet("coco.replicas")
				flags.HubXmlrpc.IsChanged = v.IsSet("hubxmlrpc.replicas")
				flags.Saline.IsChanged = v.IsSet("saline.replicas") || v.IsSet("saline.port")

				if flags.Installation.SSL.Ca.IsThirdParty() && !flags.Installation.SSL.DB.CA.IsThirdParty() {
					flags.Installation.SSL.DB.CA.Root = flags.Installation.SSL.Ca.Root
					flags.Installation.SSL.DB.CA.Intermediate = flags.Installation.SSL.Ca.Intermediate
				}
				if flags.Installation.SSL.Server.IsDefined() && !flags.Installation.SSL.DB.IsDefined() {
					flags.Installation.SSL.DB.Cert = flags.Installation.SSL.Server.Cert
					flags.Installation.SSL.DB.Key = flags.Installation.SSL.Server.Key
				}
			}
			return utils.CommandHelper(globalFlags, cmd, args, &flags, flagsUpdater, run)
		},
	}

	adm_utils.AddMirrorFlag(cmd)
	AddInstallFlags(cmd)
	podman.AddPodmanArgFlag(cmd)

	return cmd
}

// NewCommand for installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, installForPodman)
}
