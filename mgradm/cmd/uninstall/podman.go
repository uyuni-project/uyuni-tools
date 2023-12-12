// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForPodman(dryRun bool, purge bool) {

	// Uninstall the service
	podman.UninstallService("uyuni-server", dryRun)

	// Force stop the pod
	podman.DeleteContainer(podman.ServerContainerName, dryRun)

	// Remove the volumes
	if purge {
		volumes := []string{"cgroup"}
		for volume := range utils.VOLUMES {
			volumes = append(volumes, volume)
		}
		for _, volume := range volumes {
			podman.DeleteVolume(volume, dryRun)
		}
		log.Info().Msg("All volumes removed")
	}

	podman.DeleteNetwork(dryRun)

	podman.ReloadDaemon(dryRun)
}
