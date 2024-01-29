// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanInspectFlags struct {
	Image types.ImageFlags `mapstructure:",squash"`
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	podmanCmd := &cobra.Command{
		Use:   "podman",
		Short: "inspect podman image",
		Long: `Extract information from podman image and/or deployment
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanInspectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, inspectImagePodman)
		},
	}

	shared.AddInspectFlags(podmanCmd)

	return podmanCmd
}
