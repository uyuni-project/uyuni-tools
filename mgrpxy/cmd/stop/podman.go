// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package stop

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func podmanStop(
	globalFlags *types.GlobalFlags,
	flags *stopFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return systemd.StopService(podman.ProxyService)
}
