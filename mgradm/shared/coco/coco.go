// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package coco

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Upgrade coco attestation.
func Upgrade(
	systemd podman.Systemd,
	authFile string,
	cocoFlags adm_utils.CocoFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
) error {
	if cocoFlags.Image.Name == "" {
		// Don't touch the coco service in ptf if not already present.
		return nil
	}

	if err := podman.CreateCredentialsSecrets(
		podman.DBUserSecret, db.User,
		podman.DBPassSecret, db.Password,
	); err != nil {
		return err
	}

	if err := writeCocoServiceFiles(
		systemd, authFile, cocoFlags, baseImage, db,
	); err != nil {
		return err
	}

	if !cocoFlags.IsChanged {
		return systemd.RestartInstantiated(podman.ServerAttestationService)
	}
	return systemd.ScaleService(cocoFlags.Replicas, podman.ServerAttestationService)
}

func writeCocoServiceFiles(
	systemd podman.Systemd,
	authFile string,
	cocoFlags adm_utils.CocoFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
) error {
	image := cocoFlags.Image
	currentReplicas := systemd.CurrentReplicaCount(podman.ServerAttestationService)
	log.Debug().Msgf("Current Confidential Computing replicas running are %d.", currentReplicas)

	if image.Tag == "" {
		if baseImage.Tag != "" {
			image.Tag = baseImage.Tag
		} else {
			image.Tag = "latest"
		}
	}
	if !cocoFlags.IsChanged {
		log.Debug().Msg("Confidential Computing settings are not changed.")
	} else if cocoFlags.Replicas == 0 {
		log.Debug().Msg("No Confidential Computing requested.")
	}

	cocoImage, err := utils.ComputeImage(baseImage.Registry.Host, baseImage.Tag, image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	pullEnabled := (cocoFlags.Replicas > 0 && cocoFlags.IsChanged) || (currentReplicas > 0 && !cocoFlags.IsChanged)

	preparedImage, err := podman.PrepareImage(authFile, cocoImage, baseImage.PullPolicy, pullEnabled)
	if err != nil {
		return err
	}

	attestationData := templates.AttestationServiceTemplateData{
		NamePrefix:   "uyuni",
		Network:      podman.UyuniNetwork,
		Image:        preparedImage,
		DBUserSecret: podman.DBUserSecret,
		DBPassSecret: podman.DBPassSecret,
	}

	log.Info().Msg(L("Setting up confidential computing attestation service"))

	if err := utils.WriteTemplateToFile(attestationData,
		podman.GetServicePath(podman.ServerAttestationService+"@"), 0555, true); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf(`Environment=UYUNI_SERVER_ATTESTATION_IMAGE=%s
Environment=database_connection=jdbc:postgresql://%s:%d/%s
`, preparedImage, db.Host, db.Port, db.Name)

	if err := podman.GenerateSystemdConfFile(
		podman.ServerAttestationService+"@", "generated.conf", environment, true,
	); err != nil {
		return utils.Errorf(err, L("cannot generate systemd conf file"))
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}
	return nil
}

// SetupCocoContainer sets up the confidential computing attestation service.
func SetupCocoContainer(
	systemd podman.Systemd,
	authFile string,
	coco adm_utils.CocoFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
) error {
	if err := writeCocoServiceFiles(
		systemd, authFile, coco, baseImage, db,
	); err != nil {
		return err
	}
	return systemd.ScaleService(coco.Replicas, podman.ServerAttestationService)
}
