// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForPodman(
	_ *types.GlobalFlags,
	flags *utils.UninstallFlags,
	_ *cobra.Command,
	_ []string,
) error {
	// Get the images from the service configs before they are removed
	images := []string{
		podman.GetServiceImage(podman.ServerService),
		podman.GetServiceImage(podman.ServerAttestationService + "@"),
		podman.GetServiceImage(podman.HubXmlrpcService),
	}

	// Uninstall the service
	podman.UninstallService("uyuni-server", !flags.Force)
	// Force stop the pod
	podman.DeleteContainer(podman.ServerContainerName, !flags.Force)

	podman.UninstallInstantiatedService(podman.ServerAttestationService, !flags.Force)
	podman.UninstallInstantiatedService(podman.HubXmlrpcService, !flags.Force)

	// Remove the volumes
	if flags.Purge.Volumes {
		allOk := true
		volumes := []string{"cgroup"}
		for _, volume := range utils.ServerVolumeMounts {
			volumes = append(volumes, volume.Name)
		}
		for _, volume := range volumes {
			if err := podman.DeleteVolume(volume, !flags.Force); err != nil {
				log.Warn().Err(err).Msgf(L("Failed to remove volume %s"), volume)
				allOk = false
			}
		}
		if allOk {
			log.Info().Msg(L("All volumes have been removed"))
		} else {
			log.Warn().Msg(L("Some volumes have not been removed completely"))
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

	podman.DeleteNetwork(!flags.Force)

	podman.DeleteSecret(podman.DBUserSecret, !flags.Force)
	podman.DeleteSecret(podman.DBPassSecret, !flags.Force)

	err := podman.ReloadDaemon(!flags.Force)

	if !flags.Force {
		log.Warn().Msg(
			L("Nothing has been uninstalled, run with --force and --purge-volumes to actually uninstall and clear data"),
		)
	} else if !flags.Purge.Volumes {
		log.Warn().Msg(L("Data have been kept, use podman volume commands to clear the volumes"))
	}

	return err
}
