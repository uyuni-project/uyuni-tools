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

func podmanRestart(
	_ *types.GlobalFlags,
	_ *restartFlags,
	_ *cobra.Command,
	_ []string,
) error {
	err1 := podman.RestartService(podman.ServerService)
	err2 := podman.RestartInstantiated(podman.ServerAttestationService)
	err3 := podman.RestartInstantiated(podman.HubXmlrpcService)
	return utils.JoinErrors(err1, err2, err3)
}
