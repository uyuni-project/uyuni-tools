// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package restart

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

var systemd podman.Systemd = podman.NewSystemd()

func podmanRestart(
	_ *types.GlobalFlags,
	_ *restartFlags,
	_ *cobra.Command,
	_ []string,
) error {
	return errors.Join(
		systemd.ReloadDaemon(false),
		systemd.RestartService(podman.DBService),
		systemd.RestartService(podman.ServerService),
		systemd.RestartInstantiated(podman.ServerAttestationService),
		systemd.RestartInstantiated(podman.HubXmlrpcService),
	)
}
