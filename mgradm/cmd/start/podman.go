// SPDX-FileCopyrightText: 2024 SUSE LLC
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
	err1 := systemd.StartInstantiated(podman.ServerAttestationService)
	err2 := systemd.StartInstantiated(podman.HubXmlrpcService)
	err3 := systemd.StartService(podman.ServerService)
	return utils.JoinErrors(err1, err2, err3)
}
