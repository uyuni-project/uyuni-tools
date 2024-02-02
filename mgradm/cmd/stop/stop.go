// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package stop

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type stopFlags struct {
	Backend string
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "stop the server",
		Long:  "Stop the server",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags stopFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, stop)
		},
	}

	stopCmd.SetUsageTemplate(stopCmd.UsageTemplate())

	if utils.KubernetesBuilt {
		utils.AddBackendFlag(stopCmd)
	}

	return stopCmd
}

func stop(globalFlags *types.GlobalFlags, flags *stopFlags, cmd *cobra.Command, args []string) error {
	fn, err := shared.ChoosePodmanOrKubernetes(cmd.Flags(), podmanStop, kubernetesStop)
	if err != nil {
		return err
	}

	return fn(globalFlags, flags, cmd, args)
}
