// SPDX-FileCopyrightText: 2024 SUSE LLC
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

func podmanStatus(
	globalFlags *types.GlobalFlags,
	flags *statusFlags,
	cmd *cobra.Command,
	args []string,
) error {
	// Show the status and that's it if the service is not running
	if !podman.IsServiceRunning(podman.ServerService) {
		_ = utils.RunCmdStdMapping(zerolog.DebugLevel, "systemctl", "status", "--no-pager", podman.ServerService)
	} else {
		// Run spacewalk-service status in the container
		cnx := shared.NewConnection("podman", podman.ServerContainerName, "")
		_ = adm_utils.ExecCommand(zerolog.InfoLevel, cnx, "spacewalk-service", "status")
	}

	for i := 0; i < podman.CurrentReplicaCount(podman.ServerAttestationService); i++ {
		println() // add an empty line between the previous logs and this one
		_ = utils.RunCmdStdMapping(zerolog.DebugLevel, "systemctl", "status", "--no-pager", fmt.Sprintf("%s@%d", podman.ServerAttestationService, i))
	}

	for i := 0; i < podman.CurrentReplicaCount(podman.HubXmlrpcService); i++ {
		println() // add an empty line between the previous logs and this one
		_ = utils.RunCmdStdMapping(zerolog.DebugLevel, "systemctl", "status", "--no-pager", fmt.Sprintf("%s@%d", podman.HubXmlrpcService, i))
	}

	return nil
}
