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

func podmanStop(
	_ *types.GlobalFlags,
	_ *stopFlags,
	_ *cobra.Command,
	_ []string,
) error {
	err1 := podman.StopInstantiated(podman.ServerAttestationService)
	err2 := podman.StopInstantiated(podman.HubXmlrpcService)
	err3 := podman.StopService(podman.ServerService)
	return utils.JoinErrors(err1, err2, err3)
}
