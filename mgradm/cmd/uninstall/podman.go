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

func uninstallForPodman(
	globalFlags *types.GlobalFlags,
	flags *utils.UninstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	// Uninstall the service
	podman.UninstallService("uyuni-server", !flags.Force)
	// Force stop the pod
	podman.DeleteContainer(podman.ServerContainerName, !flags.Force)

	if err := podman.ScaleService(0, podman.ServerAttestationService); err != nil {
		return utils.Errorf(err, L("cannot delete confidential computing attestation instances"))
	}
	// Remove the service unit
	servicePath := podman.GetServicePath(podman.ServerAttestationService + "@")
	if _, err := os.Stat(servicePath); !os.IsNotExist(err) {
		log.Info().Msgf(L("Remove %s"), servicePath)
		if err := os.Remove(servicePath); err != nil {
			log.Error().Err(err).Msgf(L("Failed to remove %s.service file"), podman.ServerAttestationService+"@")
		}
	}

	if podman.HasService(podman.HubXmlrpcService) {
		podman.UninstallService(podman.HubXmlrpcService, !flags.Force)
		podman.DeleteContainer(podman.HubXmlrpcService, !flags.Force)
	}

	// Remove the volumes
	if flags.Purge.Volumes {
		volumes := []string{"cgroup"}
		for _, volume := range utils.ServerVolumeMounts {
			volumes = append(volumes, volume.Name)
		}
		for _, volume := range volumes {
			if err := podman.DeleteVolume(volume, !flags.Force); err != nil {
				return utils.Errorf(err, L("cannot delete volume %s"), volume)
			}
		}
		log.Info().Msg(L("All volumes removed"))
	}

	podman.DeleteNetwork(!flags.Force)

	err := podman.ReloadDaemon(!flags.Force)

	if !flags.Force {
		log.Warn().Msg(L("Nothing has been uninstalled, run with --force and --purge-volumes to actually uninstall and clear data"))
	} else if !flags.Purge.Volumes {
		log.Warn().Msg(L("Data have been kept, use podman volume commands to clear the volumes"))
	}

	return err
}
