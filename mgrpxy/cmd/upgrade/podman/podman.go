// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

// NewCommand install a new proxy on podman from scratch.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	podmanCmd := &cobra.Command{
		Use:   "podman",
		Short: L("Upgrade a proxy on podman"),
		Long: L(`Upgrade a proxy on podman

The upgrade podman command assumes podman is upgraded locally.

NOTE: for now upgrading on a remote podman is not supported!
`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podman.PodmanProxyFlags
			flags.ProxyImageFlags.Registry = globalFlags.Registry
			return shared_utils.CommandHelper(globalFlags, cmd, args, &flags, upgradePodman)
		},
	}

	utils.AddSCCFlag(podmanCmd)
	utils.AddImageFlags(podmanCmd)
	shared_podman.AddPodmanArgFlag(podmanCmd)

	return podmanCmd
}
