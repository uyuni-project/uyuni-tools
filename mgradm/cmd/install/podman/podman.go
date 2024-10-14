// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanInstallFlags struct {
	shared.InstallFlags `mapstructure:",squash"`
	Podman              podman.PodmanFlags
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[podmanInstallFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "podman [fqdn]",
		Short: L("Install a new server on podman"),
		Long: L(`Install a new server on podman

The install podman command assumes podman is installed locally.

NOTE: installing on a remote podman is not supported yet!
`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanInstallFlags
			flagsUpdater := func(v *viper.Viper) {
				flags.InstallFlags.Coco.IsChanged = v.IsSet("coco.replicas")
				flags.InstallFlags.HubXmlrpc.IsChanged = v.IsSet("hubxmlrpc.replicas")
			}
			return utils.CommandHelper(globalFlags, cmd, args, &flags, flagsUpdater, run)
		},
	}

	shared.AddInstallFlags(cmd)
	podman.AddPodmanArgFlag(cmd)

	return cmd
}

// NewCommand for podman installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, installForPodman)
}
