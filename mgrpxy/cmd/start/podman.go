// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package start

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

var systemd podman.Systemd = podman.NewSystemd()

func podmanStart(
	_ *types.GlobalFlags,
	_ *startFlags,
	_ *cobra.Command,
	_ []string,
) error {
	return systemd.StartService(podman.ProxyService)
}
