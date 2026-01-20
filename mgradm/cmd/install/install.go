// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"fmt"
	"os/exec"

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

// updateFlags handles the logic for updating flags from Viper configuration.
// Extracting this reduces the cognitive complexity of the main command function.
func updateFlags(flags *podmanInstallFlags, v *viper.Viper) {
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
	// Note: server-image and server-tag are now handled automatically by the mapstructure tags in utils.go
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
			// ensure the right amount of args, managing podman
			if len(args) > 0 && args[0] == "podman" {
				return cobra.MaximumNArgs(2)(cmd, args)
			}
			return cobra.MaximumNArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// FIX #621: Pre-flight check
			if _, err := exec.LookPath("podman"); err != nil {
				return fmt.Errorf("podman is not installed. Please install podman before running this command")
			}

			// If the alias "install podman" is used, "podman" will be the first arg.
			if len(args) > 0 && args[0] == "podman" {
				args = args[1:]
			}

			var flags podmanInstallFlags
			flagsUpdater := func(v *viper.Viper) {
				updateFlags(&flags, v)
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