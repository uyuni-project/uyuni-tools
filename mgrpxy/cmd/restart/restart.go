// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package restart

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type restartFlags struct {
	Backend string
}

// NewCommand to restart server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: L("Restart the proxy"),
		Long:  L("Restart the proxy"),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags restartFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, restart)
		},
	}
	restartCmd.SetUsageTemplate(restartCmd.UsageTemplate())

	utils.AddBackendFlag(restartCmd)

	return restartCmd
}

func restart(globalFlags *types.GlobalFlags, flags *restartFlags, cmd *cobra.Command, args []string) error {
	fn, err := shared.ChooseProxyPodmanOrKubernetes(cmd.Flags(), podmanRestart, kubernetesRestart)
	if err != nil {
		return err
	}

	return fn(globalFlags, flags, cmd, args)
}
