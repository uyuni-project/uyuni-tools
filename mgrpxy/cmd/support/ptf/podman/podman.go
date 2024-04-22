// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanPTFFlags struct {
	UpgradeFlags podman.PodmanProxyUpgradeFlags `mapstructure:",squash"`
	PTFId        string                         `mapstructure:"ptf"`
	TestId       string                         `mapstructure:"test"`
}

// NewCommand for podman installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags podmanPTFFlags
	podmanCmd := &cobra.Command{
		Use: "podman",

		Short: L("install a PTF or Test package on podman"),
		Long: L(`Install a PTF or Test package on podman

The support ptf podman command assumes podman is installed locally and
the host machine is register to SCC.

NOTE: for now installing on a remote podman is not supported!
`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return shared_utils.CommandHelper(globalFlags, cmd, args, &flags, ptfForPodman)
		},
	}

	utils.AddImageUpgradeFlags(podmanCmd)
	shared_utils.AddPTFFlag(podmanCmd)
	return podmanCmd
}
