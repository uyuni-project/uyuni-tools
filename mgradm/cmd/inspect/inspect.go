// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"

	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type inspectFlags struct {
	Image      string
	Tag        string
	PullPolicy string
}

// NewCommand for extracting information from image and deployment.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	inspectCmd := &cobra.Command{
		Use:   "inspect",
		Short: "inspect",
		Long:  "Extract information from image and deployment",
		Args:  cobra.MaximumNArgs(0),

		RunE: func(cmd *cobra.Command, args []string) error {
			var flags inspectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, inspect)
		},
	}

	inspectCmd.SetUsageTemplate(inspectCmd.UsageTemplate())
	inspectCmd.Flags().String("image", "", "Image. Leave it empty to analyze the current deployment")
	inspectCmd.Flags().String("tag", "", "Tag Image. Leave it empty to analyze the current deployment")
	utils.AddPullPolicyFlag(inspectCmd)

	if utils.KubernetesBuilt {
		utils.AddBackendFlag(inspectCmd)
	}

	return inspectCmd
}

func inspect(globalFlags *types.GlobalFlags, flags *inspectFlags, cmd *cobra.Command, args []string) error {
	fn, err := shared.ChoosePodmanOrKubernetes(cmd.Flags(), podmanInspect, kuberneteInspect)
	if err != nil {
		return err
	}
	return fn(globalFlags, flags, cmd, args)
}
