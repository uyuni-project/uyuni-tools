// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package stop

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func podmanStop(
	globalFlags *types.GlobalFlags,
	flags *stopFlags,
	cmd *cobra.Command,
	args []string,
) error {
	err1 := systemd.StopInstantiated(podman.ServerAttestationService)
	err2 := systemd.StopInstantiated(podman.HubXmlrpcService)
	err3 := systemd.StopService(podman.ServerService)
	return utils.JoinErrors(err1, err2, err3)
}
