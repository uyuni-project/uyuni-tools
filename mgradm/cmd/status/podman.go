// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func podmanStatus(
	_ *types.GlobalFlags,
	_ *statusFlags,
	_ *cobra.Command,
	_ []string,
) error {
	if systemd.HasService(podman.DBService) {
		_ = utils.RunCmdStdMapping(zerolog.DebugLevel, "systemctl", "status", "--no-pager", podman.DBService)
	}

	// Show the status and that's it if the service is not running
	if !systemd.IsServiceRunning(podman.ServerService) {
		_ = utils.RunCmdStdMapping(zerolog.DebugLevel, "systemctl", "status", "--no-pager", podman.ServerService)
	} else {
		// Run spacewalk-service status in the container
		cnx := shared.NewConnection("podman", podman.ServerContainerName, "")
		_ = adm_utils.ExecCommand(zerolog.InfoLevel, cnx, "spacewalk-service", "status")
	}

	for i := 0; i < systemd.CurrentReplicaCount(podman.ServerAttestationService); i++ {
		println() // add an empty line between the previous logs and this one
		_ = utils.RunCmdStdMapping(
			zerolog.DebugLevel, "systemctl", "status", "--no-pager", fmt.Sprintf("%s@%d", podman.ServerAttestationService, i),
		)
	}

	for i := 0; i < systemd.CurrentReplicaCount(podman.HubXmlrpcService); i++ {
		println() // add an empty line between the previous logs and this one
		_ = utils.RunCmdStdMapping(
			zerolog.DebugLevel, "systemctl", "status", "--no-pager", fmt.Sprintf("%s@%d", podman.HubXmlrpcService, i),
		)
	}

	_ = utils.RunCmdStdMapping(
		zerolog.DebugLevel, "systemctl", "status", "--no-pager", fmt.Sprintf("%s@%d", podman.EventProcessorService, 0), // refer to ScaleService // TODO: check if follow the pattern, or enforce 1 here
	)

	return nil
}
