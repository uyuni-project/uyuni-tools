// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func uninstallForPodman(
	_ *types.GlobalFlags,
	flags *utils.UninstallFlags,
	_ *cobra.Command,
	_ []string,
) error {
	dryRun := !flags.Force

	// Get the images from the service configs before they are removed
	images := []string{
		podman.GetServiceImage("uyuni-proxy-httpd"),
		podman.GetServiceImage("uyuni-proxy-salt-broker"),
		podman.GetServiceImage("uyuni-proxy-squid"),
		podman.GetServiceImage("uyuni-proxy-ssh"),
		podman.GetServiceImage("uyuni-proxy-tftpd"),
	}

	// Uninstall the service
	systemd.UninstallService("uyuni-proxy-pod", dryRun)
	systemd.UninstallService("uyuni-proxy-httpd", dryRun)
	systemd.UninstallService("uyuni-proxy-salt-broker", dryRun)
	systemd.UninstallService("uyuni-proxy-squid", dryRun)
	systemd.UninstallService("uyuni-proxy-ssh", dryRun)
	systemd.UninstallService("uyuni-proxy-tftpd", dryRun)

	// Force stop the pod
	for _, containerName := range podman.ProxyContainerNames {
		podman.DeleteContainer(containerName, dryRun)
	}

	// Remove the volumes
	if flags.Purge.Volumes {
		// Merge all proxy containers volumes into a map
		volumes := []string{}
		for _, volume := range utils.ProxyHttpdVolumes {
			volumes = append(volumes, volume.Name)
		}
		for _, volume := range utils.ProxySquidVolumes {
			volumes = append(volumes, volume.Name)
		}
		for _, volume := range utils.ProxyTftpdVolumes {
			volumes = append(volumes, volume.Name)
		}

		// Delete each volume
		for _, volume := range volumes {
			if err := podman.DeleteVolume(volume, dryRun); err != nil {
				return utils.Errorf(err, L("cannot delete volume %s"), volume)
			}
		}
		log.Info().Msg(L("All volumes removed"))
		// Remove config dir
		if err := os.RemoveAll("/etc/uyuni/proxy"); err != nil {
			log.Warn().Msg(L("Failed to delete /etc/uyuni/proxy folder"))
		} else {
			log.Info().Msg(L("/etc/uyuni/proxy folder removed"))
		}
	}

	if flags.Purge.Images {
		for _, image := range images {
			if image != "" {
				if err := podman.DeleteImage(image, !flags.Force); err != nil {
					return utils.Errorf(err, L("cannot delete image %s"), image)
				}
			}
		}
		log.Info().Msg(L("All images have been removed"))
	}

	podman.DeleteNetwork(dryRun)

	err := systemd.ReloadDaemon(dryRun)

	if dryRun {
		log.Warn().Msg(
			L("Nothing has been uninstalled, run with --force and --purge-volumes to actually uninstall and clear data"),
		)
	} else if !flags.Purge.Volumes {
		log.Warn().Msg(L("Data have been kept, use podman volume commands to clear the volumes"))
	}

	return err
}
