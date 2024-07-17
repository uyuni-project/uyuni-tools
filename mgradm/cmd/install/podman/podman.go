// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
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

// NewCommand for podman installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	podmanCmd := &cobra.Command{
		Use:   "podman [fqdn]",
		Short: L("Install a new server on podman"),
		Long: L(`Install a new server on podman

The install podman command assumes podman is installed locally.

NOTE: installing on a remote podman is not supported yet!
`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanInstallFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, installForPodman)
		},
	}

	shared.AddInstallFlags(podmanCmd)
	podman.AddPodmanArgFlag(podmanCmd)

	return podmanCmd
}
