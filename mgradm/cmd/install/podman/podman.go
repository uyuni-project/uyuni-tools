// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanInstallFlags struct {
	shared.InstallFlags `mapstructure:",squash"`
	Podman              podman.PodmanFlags
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	podmanCmd := &cobra.Command{
		Use:   "podman [fqdn]",
		Short: "install a new server on podman from scratch",
		Long: `Install a new server on podman from scratch

The install podman command assumes podman is installed locally

NOTE: for now installing on a remote podman is not supported!
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanInstallFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, installForPodman)
		},
	}

	shared.AddInstallFlags(podmanCmd)
	podman.AddPodmanInstallFlag(podmanCmd)

	return podmanCmd
}
