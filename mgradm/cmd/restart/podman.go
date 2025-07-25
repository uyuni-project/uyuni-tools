// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package restart

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func podmanRestart(
	_ *types.GlobalFlags,
	_ *restartFlags,
	_ *cobra.Command,
	_ []string,
) error {
	return utils.JoinErrors(
		systemd.RestartService(podman.DBService),
		systemd.RestartService(podman.ServerService),
		systemd.RestartInstantiated(podman.ServerAttestationService),
		systemd.RestartInstantiated(podman.HubXmlrpcService),
		systemd.RestartInstantiated(podman.EventProcessorService),
	)
}
