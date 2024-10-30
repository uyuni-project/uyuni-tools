// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package start

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func podmanStart(
	_ *types.GlobalFlags,
	_ *startFlags,
	_ *cobra.Command,
	_ []string,
) error {
	return utils.JoinErrors(
		systemd.StartService(podman.DBService),
		systemd.StartInstantiated(podman.ServerAttestationService),
		systemd.StartInstantiated(podman.HubXmlrpcService),
		systemd.StartService(podman.ServerService),
	)
}
