// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package start

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type startFlags struct {
	Backend string
}

// NewCommand starts the server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start the proxy",
		Long:  "Start the proxy",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags startFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, start)
		},
	}
	startCmd.SetUsageTemplate(startCmd.UsageTemplate())

	if utils.KubernetesBuilt {
		utils.AddBackendFlag(startCmd)
	}

	return startCmd
}

func start(globalFlags *types.GlobalFlags, flags *startFlags, cmd *cobra.Command, args []string) error {
	fn, err := shared.ChooseProxyPodmanOrKubernetes(cmd.Flags(), podmanStart, kubernetesStart)
	if err != nil {
		return err
	}

	return fn(globalFlags, flags, cmd, args)
}
