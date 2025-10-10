// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.NewSystemd()

func podmanCacheClear(
	_ *types.GlobalFlags,
	_ *cacheClearFlags,
	_ *cobra.Command,
	_ []string,
) error {
	cnx := shared.NewConnection("podman", "uyuni-proxy-squid", "")

	if _, err := cnx.Exec("sh", "-c", "rm -rf /var/cache/squid/*"); err != nil {
		return utils.Errorf(err, L("failed to remove cached data"))
	}

	return systemd.RestartService(podman.ProxyService)
}
