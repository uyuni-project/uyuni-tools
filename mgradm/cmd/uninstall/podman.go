// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
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
	podman.UninstallService("uyuni-server", flags.DryRun)

	// Force stop the pod
	podman.DeleteContainer(podman.ServerContainerName, flags.DryRun)

	// Remove the volumes
	if flags.PurgeVolumes {
		volumes := []string{"cgroup"}
		for volume := range utils.VOLUMES {
			volumes = append(volumes, volume)
		}
		for _, volume := range volumes {
			podman.DeleteVolume(volume, flags.DryRun)
		}
		log.Info().Msg("All volumes removed")
	}

	podman.DeleteNetwork(flags.DryRun)

	podman.ReloadDaemon(flags.DryRun)

	return nil
}
