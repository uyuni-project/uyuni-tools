// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

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
	UpgradeFlags podman.PodmanProxyFlags `mapstructure:",squash"`
	SCC          types.SCCCredentials    `mapstructure:"scc"`
	PTFId        string                  `mapstructure:"ptf"`
	TestID       string                  `mapstructure:"test"`
	CustomerID   string                  `mapstructure:"user"`
}

func newCmd(globalFlags *types.GlobalFlags, run shared_utils.CommandFunc[podmanPTFFlags]) *cobra.Command {
	var flags podmanPTFFlags
	podmanCmd := &cobra.Command{
		Use: "podman",

		Short: L("Install a PTF or Test package on podman"),
		Long: L(`Install a PTF or Test package on podman

The support ptf podman command assumes podman is installed locally and
the host machine is registered to SCC.

NOTE: for now installing on a remote podman is not supported!
`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return shared_utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	utils.AddSCCFlag(podmanCmd)
	utils.AddImageFlags(podmanCmd)
	shared_utils.AddPTFFlag(podmanCmd)
	return podmanCmd
}

// NewCommand for podman installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, ptfForPodman)
}
