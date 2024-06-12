// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"os"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
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
		volumes := []string{}
		for _, volume := range utils.PROXY_HTTPD_VOLUMES {
			volumes = append(volumes, volume.Name)
		}
		for _, volume := range utils.PROXY_SQUID_VOLUMES {
			volumes = append(volumes, volume.Name)
		}
		for _, volume := range utils.PROXY_TFTPD_VOLUMES {
			volumes = append(volumes, volume.Name)
		}

		// Delete each volume
		for _, volume := range volumes {
			if err := podman.DeleteVolume(volume, dryRun); err != nil {
				return utils.Errorf(err, L("cannot delete volume %s"), volume)
			}
		}
		log.Info().Msg(L("All volumes removed"))
		//Remove config dir
		if err := os.RemoveAll("/etc/uyuni/proxy"); err != nil {
			log.Warn().Msg(L("Failed to delete /etc/uyuni/proxy folder"))
		} else {
			log.Info().Msg(L("/etc/uyuni/proxy folder removed"))
		}
	}

	podman.DeleteNetwork(dryRun)

	err := podman.ReloadDaemon(dryRun)

	if dryRun {
		log.Warn().Msg(L("Nothing has been uninstalled, run with --force and --purgeVolumes to actually uninstall and clear data"))
	} else if !purge {
		log.Warn().Msg(L("Data have been kept, use podman volume commands to clear the volumes"))
	}

	return err
}
