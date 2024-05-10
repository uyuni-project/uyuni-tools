// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package restart

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func podmanRestart(
	globalFlags *types.GlobalFlags,
	flags *restartFlags,
	cmd *cobra.Command,
	args []string,
) error {
	err := podman.RestartService(podman.ServerService)
	if err != nil {
		return err
	}
	if podman.HasService(podman.ServerAttestationService) {
		return podman.RestartService(podman.ServerAttestationService)
	}
	if podman.HasService(podman.HubXmlrpcService) {
		return podman.RestartService(podman.HubXmlrpcService)
	}
	return nil
}
