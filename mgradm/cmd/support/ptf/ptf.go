// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package ptf

import (
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanPTFFlags struct {
	adm_utils.ServerFlags `mapstructure:",squash"`
	Podman                podman.PodmanFlags
	PTFId                 string               `mapstructure:"ptf"`
	TestID                string               `mapstructure:"test"`
	CustomerID            string               `mapstructure:"user"`
	SCC                   types.SCCCredentials `mapstructure:"scc"`
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[podmanPTFFlags]) *cobra.Command {
	podmanCmd := &cobra.Command{
		Use:     "ptf",
		Aliases: []string{"ptf podman"},
		Short:   L("Install a PTF or Test package on podman"),
		Long: L(`Install a PTF or Test package on podman

The support ptf podman command assumes podman is installed locally and
the host machine is register to SCC.

NOTE: for now installing on a remote podman is not supported!
`),
		Args: func(cmd *cobra.Command, args []string) error {
			// ensure the right amount of args, managing podman
			if len(args) > 0 && args[0] == "podman" {
				return cobra.MaximumNArgs(1)(cmd, args)
			}
			return cobra.MaximumNArgs(0)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanPTFFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	adm_utils.AddSCCFlag(podmanCmd)
	utils.AddPTFFlag(podmanCmd)
	utils.AddPullPolicyFlag(podmanCmd)
	utils.AddRegistryFlag(podmanCmd)

	return podmanCmd
}

// NewCommand for podman installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, ptfForPodman)
}
