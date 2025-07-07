// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package saltEventProcessor

import (
	"github.com/rs/zerolog/log"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Upgrade salt event processor
func Upgrade(
	systemd podman.Systemd,
	authFile string,
	registry string,
	saltEventProcessorFlags adm_utils.SaltEventProcessorFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
) error {
	if saltEventProcessorFlags.Image.Name == "" {
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

func writeSaltEventProcessorFiles(
	systemd podman.Systemd,
	authFile string,
	registry string,
	saltEventProcessorFlags adm_utils.SaltEventProcessorFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
) error {
	image := saltEventProcessorFlags.Image
	currentReplicas := systemd.CurrentReplicaCount(podman.SaltEventProcessorService)
	log.Debug().Msgf("Current running Salt event processor replicas are %d", currentReplicas)

	if image.Tag == "" {
		if baseImage.Tag == "" {
			image.Tag = baseImage.Tag
		} else {
			image.Tag = "latest"
		}
	}

	if !saltEventProcessorFlags.IsChanged {
		log.Debug().Msgf("Salt event processor settings are not changed.")
	}

	if saltEventProcessorFlags.Replicas == 0 {
		log.Debug().Msgf("No Salt event processor server requested")
	}

	saltEventProcessorImage, err := utils.ComputeImage(registry, baseImage.Tag, image)
	if err != nil {
		return utils.Errorf(err, L("Failed to compute salt event processor image URL"))
	}

	pullEnabled := (saltEventProcessorFlags.IsChanged && saltEventProcessorFlags.Replicas > 0) ||
		(!saltEventProcessorFlags.IsChanged && currentReplicas > 0)

	preparedImage, err := podman.PrepareImage(authFile, saltEventProcessorImage, baseImage.PullPolicy, pullEnabled)
	if err != nil {
		return err
	}
}

func SetupSaltEventProcessorContainer(
	systemd podman.Systemd,
	authFile string,
	registry string,
	saltEventProcessorFlags adm_utils.SaltEventProcessorFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
) error {
	if err := writeSaltEventProcessorFiles(
		systemd, authFile, registry, saltEventProcessorFlags, baseImage, db,
	); err != nil {
		return err
	}
	return systemd.ScaleService(saltEventProcessorFlags.Replicas, podman.SaltEventProcessorService)
}
