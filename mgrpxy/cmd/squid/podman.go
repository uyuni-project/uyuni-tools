// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package squid

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func podmanSquidClear(
	globalFlags *types.GlobalFlags,
	flags *squidClearFlags,
	cmd *cobra.Command,
	args []string,
) error {
	volumeName := "uyuni-proxy-squid-cache"

	if err := podman.StopService(podman.ProxyService); err != nil {
		return err
	}

	if err := podman.DeleteVolume(volumeName, false); err != nil {
		return err
	}

	if err := podman.CreateVolume(volumeName, false); err != nil {
		return err
	}

	return podman.StartService(podman.ProxyService)
}
