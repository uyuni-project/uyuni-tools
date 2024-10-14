// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package start

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type startFlags struct {
	Backend string
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[startFlags]) *cobra.Command {
	startCmd := &cobra.Command{
		Use:     "start",
		GroupID: "management",
		Short:   L("Start the server"),
		Long:    L("Start the server"),
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags startFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, run)
		},
	}
	startCmd.SetUsageTemplate(startCmd.UsageTemplate())

	if utils.KubernetesBuilt {
		utils.AddBackendFlag(startCmd)
	}

	return startCmd
}

// NewCommand starts the server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, start)
}

func start(globalFlags *types.GlobalFlags, flags *startFlags, cmd *cobra.Command, args []string) error {
	fn, err := shared.ChoosePodmanOrKubernetes(cmd.Flags(), podmanStart, kubernetesStart)
	if err != nil {
		return err
	}

	return fn(globalFlags, flags, cmd, args)
}
