// SPDX-FileCopyrightText: 2024 SUSE LLC
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

var systemd podman.Systemd = podman.SystemdImpl{}

func podmanCacheClear(
	globalFlags *types.GlobalFlags,
	flags *cacheClearFlags,
	cmd *cobra.Command,
	args []string,
) error {
	cnx := shared.NewConnection("podman", "uyuni-proxy-squid", "")

	if _, err := cnx.Exec("sh", "-c", "rm -rf /var/cache/squid/*"); err != nil {
		return utils.Errorf(err, L("failed to remove cached data"))
	}

	if _, err := cnx.Exec("sh", "-c", "squid -z --foreground"); err != nil {
		return utils.Errorf(err, L("failed to re-create the cache directories"))
	}

	return systemd.RestartService(podman.ProxyService)
}
