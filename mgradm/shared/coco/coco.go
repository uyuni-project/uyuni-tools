// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package coco

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Start starts all coco replicas.
func Start() error {
	for i := 0; i < podman.CurrentReplicaCount(podman.ServerAttestationService); i++ {
		if err := podman.StartService(fmt.Sprintf("%s@%d", podman.ServerAttestationService, i)); err != nil {
			return err
		}
	}
	return nil
}

// Stop stops all coco replicas.
func Stop() error {
	for i := 0; i < podman.CurrentReplicaCount(podman.ServerAttestationService); i++ {
		if err := podman.StopService(fmt.Sprintf("%s@%d", podman.ServerAttestationService, i)); err != nil {
			return err
		}
	}
	return nil
}

// Upgrade coco attestation.
func Upgrade(image types.ImageFlags, baseImage types.ImageFlags, dbPort int, dbName string, dbUser string, dbPassword string) error {
	if err := Stop(); err != nil {
		return err
	}
	if err := writeCocoServiceFiles(image, baseImage, dbName, dbPort, dbUser, dbPassword); err != nil {
		return err
	}
	return Start()
}

// Uninstall scales coco service to 0 and removes service files.
func Uninstall(dryRun bool) error {
	if dryRun {
		log.Info().Msg(L("Would remove uyuni-server-attestation instances."))
	} else {
		if err := podman.ScaleService(0, podman.ServerAttestationService); err != nil {
			return utils.Errorf(err, L("cannot delete confidential computing attestation instances"))
		}
	}

	name := podman.ServerAttestationService + "@"
	servicePath := podman.GetServicePath(name)
	if dryRun {
		log.Info().Msgf(L("Would remove %s"), servicePath)
	} else {
		// Remove the service unit
		if _, err := os.Stat(servicePath); !os.IsNotExist(err) {
			log.Info().Msgf(L("Remove %s"), servicePath)
			if err := os.Remove(servicePath); err != nil {
				log.Error().Err(err).Msgf(L("Failed to remove %s.service file"), podman.ServerAttestationService+"@")
			}
		}
	}

	serviceConfFolder := podman.GetServiceConfFolder(name)
	serviceConfPath := podman.GetServiceConfPath(name)
	if utils.FileExists(serviceConfFolder) {
		if utils.FileExists(serviceConfPath) {
			if dryRun {
				log.Info().Msgf(L("Would remove %s"), serviceConfPath)
			} else {
				log.Info().Msgf(L("Remove %s"), serviceConfPath)
				if err := os.Remove(serviceConfPath); err != nil {
					log.Error().Err(err).Msgf(L("Failed to remove %s file"), serviceConfPath)
				}
			}
		}

		if dryRun {
			log.Info().Msgf(L("Would remove %s if empty"), serviceConfFolder)
		} else {
			if utils.IsEmptyDirectory(serviceConfFolder) {
				log.Debug().Msgf("Removing %s folder, since it's empty", serviceConfFolder)
				_ = utils.RemoveDirectory(serviceConfFolder)
			} else {
				log.Warn().Msgf(L("%s folder contains file created by the user. Please remove them when uninstallation is completed."), serviceConfFolder)
			}
		}
	}
	return nil
}

func writeCocoServiceFiles(image types.ImageFlags, baseImage types.ImageFlags, dbName string, dbPort int, dbUser string, dbPassword string) error {
	if image.Tag == "" {
		if baseImage.Tag != "" {
			image.Tag = baseImage.Tag
		} else {
			image.Tag = "latest"
		}
	}
	cocoImage, err := utils.ComputeImage(image)
	if err != nil {
		baseImage.Tag = image.Tag
		cocoImage, err = utils.ComputeImage(baseImage, "-attestation")
		if err != nil {
			return utils.Errorf(err, L("failed to compute image URL"))
		}
	}

	attestationData := templates.AttestationServiceTemplateData{
		NamePrefix: "uyuni",
		Network:    podman.UyuniNetwork,
		Image:      cocoImage,
	}

	log.Info().Msg(L("Setting up confidential computing attestation service"))

	if err := utils.WriteTemplateToFile(attestationData,
		podman.GetServicePath(podman.ServerAttestationService+"@"), 0555, true); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf(`Environment=UYUNI_IMAGE=%s
	Environment=database_connection=jdbc:postgresql://uyuni-server.mgr.internal:%d/%s
	Environment=database_user=%s
	Environment=database_password=%s`, cocoImage, dbPort, dbName, dbUser, dbPassword)

	if err := podman.GenerateSystemdConfFile(podman.ServerAttestationService+"@", "Service", environment); err != nil {
		return utils.Errorf(err, L("cannot generate systemd conf file"))
	}

	if err := podman.ReloadDaemon(false); err != nil {
		return err
	}
	return nil
}

// SetupCocoContainer sets up the confidential computing attestation service.
func SetupCocoContainer(replicas int, image types.ImageFlags, baseImage types.ImageFlags, dbName string, dbPort int, dbUser string, dbPassword string) error {
	if err := writeCocoServiceFiles(image, baseImage, dbName, dbPort, dbUser, dbPassword); err != nil {
		return err
	}
	return podman.ScaleService(replicas, podman.ServerAttestationService)
}
