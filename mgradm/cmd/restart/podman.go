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
		if err := podman.RestartService(podman.ServerAttestationService); err != nil {
			return err
		}
	}
	if podman.HasService(podman.HubXmlrpcService) {
		if err := podman.RestartService(podman.HubXmlrpcService); err != nil {
			return err
		}
	}
	return nil
}
