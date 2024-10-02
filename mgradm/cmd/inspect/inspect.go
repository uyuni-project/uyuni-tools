// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"github.com/spf13/cobra"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// InspectFlags are the flags used by inspect commands.
type inspectFlags struct {
	Image   types.ImageFlags `mapstructure:",squash"`
	SCC     types.SCCCredentials
	Backend string
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[inspectFlags]) *cobra.Command {
	inspectCmd := &cobra.Command{
		Use:     "inspect",
		GroupID: "deploy",
		Short:   L("Inspect"),
		Long:    L("Extract information from image and deployment"),
		Args:    cobra.MaximumNArgs(0),

		RunE: func(cmd *cobra.Command, args []string) error {
			var flags inspectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, run)
		},
	}

	inspectCmd.SetUsageTemplate(inspectCmd.UsageTemplate())

	cmd_utils.AddSCCFlag(inspectCmd)
	cmd_utils.AddImageFlag(inspectCmd)

	if utils.KubernetesBuilt {
		utils.AddBackendFlag(inspectCmd)
	}

	return inspectCmd
}

// NewCommand for extracting information from image and deployment.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, inspect)
}

func inspect(globalFlags *types.GlobalFlags, flags *inspectFlags, cmd *cobra.Command, args []string) error {
	fn, err := shared.ChoosePodmanOrKubernetes(cmd.Flags(), podmanInspect, kuberneteInspect)
	if err != nil {
		return err
	}
	return fn(globalFlags, flags, cmd, args)
}
