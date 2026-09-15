// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package disable

import (
	"github.com/spf13/cobra"
	admPodman "github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.NewSystemd()

func podmanDisable(
	_ *types.GlobalFlags,
	_ *disableFlags,
	_ *cobra.Command,
	_ []string,
) error {
	return admPodman.DisableServices(systemd)
}

type disableFlags struct {
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[disableFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable",
		GroupID: "management",
		Short:   L("Disable server services at boot"),
		Long:    L("Disable server services at boot without stopping them"),
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags disableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	cmd.SetUsageTemplate(cmd.UsageTemplate())
	return cmd
}

// NewCommand disables automatic startup of the server services.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, podmanDisable)
}
