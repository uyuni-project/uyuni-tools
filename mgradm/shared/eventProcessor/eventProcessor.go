// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package eventProcessor

import (
	"github.com/rs/zerolog/log"
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
	// call podman secret create to store the secret from temp file to podman
	if err := podman.CreateCredentialsSecrets(
		podman.DBUserSecret, db.User,
		podman.DBPassSecret, db.Password,
	); err != nil {
		return err
	}

}

func writeEventProcessorFiles(
	systemd podman.Systemd,
	authFile string,
	registry string,
	eventProcessorFlags adm_utils.EventProcessorFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
) error {
	image := eventProcessorFlags.Image
	currentReplicas := systemd.CurrentReplicaCount(podman.EventProcessorService)
	log.Debug().Msgf("Current running Salt event processor replicas are %d", currentReplicas)

	if image.Tag == "" {
		if baseImage.Tag == "" {
			image.Tag = baseImage.Tag
		} else {
			image.Tag = "latest"
		}
	}

	if !eventProcessorFlags.IsChanged {
		log.Debug().Msgf("Salt event processor settings are not changed.")
	}

	if eventProcessorFlags.Replicas == 0 {
		log.Debug().Msgf("No Salt event processor server requested")
	}

	saltEventProcessorImage, err := utils.ComputeImage(registry, baseImage.Tag, image)
	if err != nil {
		return utils.Errorf(err, L("Failed to compute salt event processor image URL"))
	}

	pullEnabled := (eventProcessorFlags.IsChanged && eventProcessorFlags.Replicas > 0) ||
		(!eventProcessorFlags.IsChanged && currentReplicas > 0)

	preparedImage, err := podman.PrepareImage(authFile, saltEventProcessorImage, baseImage.PullPolicy, pullEnabled)
	if err != nil {
		return err
	}
}

func SetupEventProcessorContainer(
	systemd podman.Systemd,
	authFile string,
	registry string,
	eventProcessorFlags adm_utils.EventProcessorFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
) error {
	if err := writeEventProcessorFiles(
		systemd, authFile, registry, eventProcessorFlags, baseImage, db,
	); err != nil {
		return err
	}
	return systemd.ScaleService(eventProcessorFlags.Replicas, podman.EventProcessorService)
}
