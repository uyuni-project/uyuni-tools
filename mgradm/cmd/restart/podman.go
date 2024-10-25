// SPDX-FileCopyrightText: 2024 SUSE LLC
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
	globalFlags *types.GlobalFlags,
	flags *restartFlags,
	cmd *cobra.Command,
	args []string,
) error {
	err1 := systemd.RestartService(podman.ServerService)
	err2 := systemd.RestartInstantiated(podman.ServerAttestationService)
	err3 := systemd.RestartInstantiated(podman.HubXmlrpcService)
	return utils.JoinErrors(err1, err2, err3)
}
