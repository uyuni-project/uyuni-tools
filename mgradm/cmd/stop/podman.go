// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package stop

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func podmanStop(
	globalFlags *types.GlobalFlags,
	flags *stopFlags,
	cmd *cobra.Command,
	args []string,
) error {
	for i := 0; i < podman.CurrentReplicaCount(podman.ServerAttestationService); i++ {
		if err := podman.StopService(fmt.Sprintf("%s@%d", podman.ServerAttestationService, i)); err != nil {
			return err
		}
	}
	return podman.StopService(podman.ServerService)
}
