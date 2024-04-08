// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package start

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func podmanStart(
	globalFlags *types.GlobalFlags,
	flags *startFlags,
	cmd *cobra.Command,
	args []string,
) error {
	if podman.HasService(podman.ServerAttestationService) {
		if err := podman.StartService(podman.ServerAttestationService); err != nil {
			return err
		}
	}
	return podman.StartService(podman.ServerService)
}
