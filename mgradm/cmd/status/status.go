// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"errors"

	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type statusFlags struct {
}

// NewCommand to get the status of the server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: L("Get the server status"),
		Long:  L("Get the server status"),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags statusFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, status)
		},
	}
	cmd.SetUsageTemplate(cmd.UsageTemplate())

	return cmd
}

func status(globalFlags *types.GlobalFlags, flags *statusFlags, cmd *cobra.Command, args []string) error {
	if podman.HasService(podman.ServerService) {
		return podmanStatus(globalFlags, flags, cmd, args)
	}

	if utils.IsInstalled("kubectl") && utils.IsInstalled("helm") {
		return kubernetesStatus(globalFlags, flags, cmd, args)
	}

	return errors.New(L("no installed server detected"))
}
