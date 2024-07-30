// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"github.com/spf13/cobra"
	inspect_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// NewCommand for extracting information from image and deployment.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	inspectCmd := &cobra.Command{
		Use:     "inspect",
		GroupID: "deploy",
		Short:   L("Inspect"),
		Long:    L("Extract information from image and deployment"),
		Args:    cobra.MaximumNArgs(0),

		RunE: func(cmd *cobra.Command, args []string) error {
			var flags inspect_shared.InspectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, inspect)
		},
	}

	inspectCmd.SetUsageTemplate(inspectCmd.UsageTemplate())

	inspect_shared.AddInspectFlags(inspectCmd)

	if utils.KubernetesBuilt {
		utils.AddBackendFlag(inspectCmd)
	}

	return inspectCmd
}

func inspect(globalFlags *types.GlobalFlags, flags *inspect_shared.InspectFlags, cmd *cobra.Command, args []string) error {
	fn, err := shared.ChoosePodmanOrKubernetes(cmd.Flags(), podmanInspect, kuberneteInspect)
	if err != nil {
		return err
	}
	return fn(globalFlags, flags, cmd, args)
}
