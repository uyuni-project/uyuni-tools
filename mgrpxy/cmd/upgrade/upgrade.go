// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package upgrade

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

func newCmd(globalFlags *types.GlobalFlags, run shared_utils.CommandFunc[podman.PodmanProxyFlags]) *cobra.Command {
	podmanCmd := &cobra.Command{
		Use:     "upgrade",
		Aliases: []string{"upgrade podman"},
		GroupID: "deploy",
		Short:   L("Upgrade a proxy on podman"),
		Long: L(`Upgrade a proxy on podman

The upgrade podman command assumes podman is upgraded locally.

/etc/uyuni/proxy/apache.conf and /etc/uyuni/squid.conf will be used as tuning files
for apache and squid if available and not superseded by the matching command arguments.

NOTE: for now upgrading on a remote podman is not supported!
`),
		Args: func(cmd *cobra.Command, args []string) error {
			// ensure the right amount of args, managing podman
			if len(args) > 0 && args[0] == "podman" {
				return cobra.ExactArgs(1)(cmd, args)
			}
			return cobra.ExactArgs(0)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podman.PodmanProxyFlags
			return shared_utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	utils.AddSCCFlag(podmanCmd)
	utils.AddImageFlags(podmanCmd)
	shared_podman.AddPodmanArgFlag(podmanCmd)

	return podmanCmd
}

// NewCommand install a new proxy on podman from scratch.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, upgradePodman)
}
