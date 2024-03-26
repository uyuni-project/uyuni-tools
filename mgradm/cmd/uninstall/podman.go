// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForPodman(
	globalFlags *types.GlobalFlags,
	flags *uninstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	// Uninstall the service
	podman.UninstallService("uyuni-server", !flags.Force)

	// Force stop the pod
	podman.DeleteContainer(podman.ServerContainerName, !flags.Force)

	// Remove the volumes
	if flags.PurgeVolumes {
		volumes := []string{"cgroup"}
		for _, volume := range utils.ServerVolumeMounts {
			volumes = append(volumes, volume.Name)
		}
		for _, volume := range volumes {
			if err := podman.DeleteVolume(volume, !flags.Force); err != nil {
				return fmt.Errorf("cannot delete volume %s: %s", volume, err)
			}
		}
		log.Info().Msg("All volumes removed")
	}

	podman.DeleteNetwork(!flags.Force)

	return podman.ReloadDaemon(!flags.Force)
}
