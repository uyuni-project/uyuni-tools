// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package enable

import (
	"github.com/spf13/cobra"
	admPodman "github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.NewSystemd()

func podmanEnable(
	_ *types.GlobalFlags,
	_ *enableFlags,
	_ *cobra.Command,
	_ []string,
) error {
	return admPodman.EnableServices(systemd)
}

type enableFlags struct {
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[enableFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "enable",
		GroupID: "management",
		Short:   L("Enable server services at boot"),
		Long:    L("Enable server services at boot without starting them"),
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags enableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	cmd.SetUsageTemplate(cmd.UsageTemplate())
	return cmd
}

// NewCommand enables automatic startup of the server services.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, podmanEnable)
}
