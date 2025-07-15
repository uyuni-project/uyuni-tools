// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package eventProcessor

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Upgrade event processor
func Upgrade(
	systemd podman.Systemd,
	authFile string,
	registry string,
	eventProcessorFlags adm_utils.EventProcessorFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
) error {
	if eventProcessorFlags.Image.Name == "" {
		return nil
	}

	if err := writeEventProcessorFiles(
		systemd, authFile, registry, eventProcessorFlags, baseImage,
	); err != nil {
		return err
	}

	if !eventProcessorFlags.IsChanged {
		return systemd.RestartInstantiated(podman.EventProcessorService)
	}
	return systemd.ScaleService(1, podman.EventProcessorService) // TODO: we can't scale here with 1 replica, what to upgrade?
}

func writeEventProcessorFiles(
	systemd podman.Systemd,
	authFile string,
	registry string,
	eventProcessorFlags adm_utils.EventProcessorFlags,
	baseImage types.ImageFlags,
) error {
	image := eventProcessorFlags.Image

	log.Debug().Msgf("Current running event processor replica is enforced to be 1")

	if image.Tag == "" {
		if baseImage.Tag == "" {
			image.Tag = baseImage.Tag
		} else {
			image.Tag = "latest"
		}
	}

	if !eventProcessorFlags.IsChanged {
		log.Debug().Msgf("Event processor settings are not changed.")
	}

	eventProcessorImage, err := utils.ComputeImage(registry, baseImage.Tag, image)
	if err != nil {
		return utils.Errorf(err, L("Failed to compute event processor image URL"))
	}

	// Always enable pulling if service is requested (since we enforced single replica)
	preparedImage, err := podman.PrepareImage(authFile, eventProcessorImage, baseImage.PullPolicy, true)
	if err != nil {
		return err
	}

	eventProcessorData := templates.EventProcessorServiceTemplateData{
		NamePrefix:   "uyuni",
		Network:      podman.UyuniNetwork,
		DBUserSecret: podman.DBUserSecret,
		DBPassSecret: podman.DBPassSecret,
		DBName:       "susemanager", // TODO: check if we should hard code it here or set in the systemd config file
		DBPort:       utils.DBPorts,
		DBHost:       "db",
	}

	log.Info().Msg(L("Setting up event processor service"))

	if err := utils.WriteTemplateToFile(
		eventProcessorData, podman.GetServicePath(podman.EventProcessorService+"@"), 0555, true,
	); err != nil {
		return utils.Errorf(err, L("Failed to generate systemd service unit file"))
	}

	// TODO: check if we should code DB related environment in systemd conf
	environment := fmt.Sprintf(`Environment=UYUNI_EVENT_PROCESSOR_IMAGE=%s`, preparedImage) // UYUNI_EVENT_PROCESSOR_IMAGE is used in template

	if err := podman.GenerateSystemdConfFile(
		podman.EventProcessorService+"@", "generated.conf", environment, true,
	); err != nil {
		return utils.Errorf(err, L("cannot generate systemd user configuration file"))
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}

	return nil
}

func SetupEventProcessorContainer(
	systemd podman.Systemd,
	authFile string,
	registry string,
	eventProcessorFlags adm_utils.EventProcessorFlags,
	baseImage types.ImageFlags,
) error {
	if err := writeEventProcessorFiles(
		systemd, authFile, registry, eventProcessorFlags, baseImage,
	); err != nil {
		return err
	}
	// Enforce one replica
	return systemd.ScaleService(1, podman.EventProcessorService)
}
