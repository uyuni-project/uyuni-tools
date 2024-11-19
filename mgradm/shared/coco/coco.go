// SPDX-FileCopyrightText: 2024 SUSE LLC
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
	registry string,
	cocoFlags adm_utils.CocoFlags,
	baseImage types.ImageFlags,
	dbPort int,
	dbName string,
	dbUser string,
	dbPassword string,
) error {
	if cocoFlags.Image.Name == "" {
		// Don't touch the coco service in ptf if not already present.
		return nil
	}

	if err := podman.CreateDBSecrets(dbUser, dbPassword); err != nil {
		return err
	}

	if err := writeCocoServiceFiles(
		systemd, authFile, registry, cocoFlags, baseImage, dbName, dbPort,
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
	registry string,
	cocoFlags adm_utils.CocoFlags,
	baseImage types.ImageFlags,
	dbName string,
	dbPort int,
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

	cocoImage, err := utils.ComputeImage(registry, baseImage.Tag, image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	pullEnabled := (cocoFlags.Replicas > 0 && cocoFlags.IsChanged) || (currentReplicas > 0 && !cocoFlags.IsChanged)

	preparedImage, err := podman.PrepareImage(authFile, cocoImage, baseImage.PullPolicy, pullEnabled)
	if err != nil {
		return err
	}

	attestationData := templates.AttestationServiceTemplateData{
		NamePrefix: "uyuni",
		Network:    podman.UyuniNetwork,
		Image:      preparedImage,
	}

	log.Info().Msg(L("Setting up confidential computing attestation service"))

	if err := utils.WriteTemplateToFile(attestationData,
		podman.GetServicePath(podman.ServerAttestationService+"@"), 0555, true); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf(`Environment=UYUNI_IMAGE=%s
Environment=database_connection=jdbc:postgresql://uyuni-server.mgr.internal:%d/%s
`, preparedImage, dbPort, dbName)

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
	registry string,
	coco adm_utils.CocoFlags,
	baseImage types.ImageFlags,
	dbName string,
	dbPort int,
) error {
	if err := writeCocoServiceFiles(
		systemd, authFile, registry, coco, baseImage, dbName, dbPort,
	); err != nil {
		return err
	}
	return systemd.ScaleService(coco.Replicas, podman.ServerAttestationService)
}
