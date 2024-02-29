// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForPodman(dryRun bool, purge bool) error {
	// Uninstall the service
	podman.UninstallService("uyuni-proxy-pod", dryRun)
	podman.UninstallService("uyuni-proxy-httpd", dryRun)
	podman.UninstallService("uyuni-proxy-salt-broker", dryRun)
	podman.UninstallService("uyuni-proxy-squid", dryRun)
	podman.UninstallService("uyuni-proxy-ssh", dryRun)
	podman.UninstallService("uyuni-proxy-tftpd", dryRun)

	// Force stop the pod
	for _, containerName := range podman.ProxyContainerNames {
		podman.DeleteContainer(containerName, dryRun)
	}

	// Remove the volumes
	if purge {
		// Merge all proxy containers volumes into a map
		volumes := map[string]string{}
		allProxyVolumes := []map[string]string{
			utils.PROXY_HTTPD_VOLUMES,
			utils.PROXY_SQUID_VOLUMES,
			utils.PROXY_TFTPD_VOLUMES,
		}
		for _, volumesList := range allProxyVolumes {
			for volume, mount := range volumesList {
				volumes[volume] = mount
			}
		}

		// Delete each volume
		for volume := range volumes {
			if err := podman.DeleteVolume(volume, dryRun); err != nil {
				return fmt.Errorf("cannot delete volume %s: %s", volume, err)
			}
		}
		log.Info().Msg("All volumes removed")
	}

	podman.DeleteNetwork(dryRun)

	return podman.ReloadDaemon(dryRun)
}
