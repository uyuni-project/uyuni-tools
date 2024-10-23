// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package podman

import (
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanPTFFlags struct {
	Image      types.ImageFlags `mapstructure:",squash"`
	PTFId      string           `mapstructure:"ptf"`
	TestId     string           `mapstructure:"test"`
	CustomerId string           `mapstructure:"user"`
	SCC        types.SCCCredentials
	Coco       adm_utils.CocoFlags
	Hubxmlrpc  adm_utils.HubXmlrpcFlags
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[podmanPTFFlags]) *cobra.Command {
	podmanCmd := &cobra.Command{
		Use: "podman",

		Short: L("Install a PTF or Test package on podman"),
		Long: L(`Install a PTF or Test package on podman

The support ptf podman command assumes podman is installed locally and
the host machine is register to SCC.

NOTE: for now installing on a remote podman is not supported!
`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanPTFFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	utils.AddPTFFlag(podmanCmd)

	return podmanCmd
}

// NewCommand for podman installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, ptfForPodman)
}
