// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

var systemd shared_podman.Systemd = shared_podman.SystemdImpl{}

func upgradePodman(
	globalFlags *types.GlobalFlags,
	flags *podman.PodmanProxyFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return podman.Upgrade(systemd, globalFlags, flags, cmd, args)
}
