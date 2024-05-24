// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package coco

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// SetupCocoContainer sets up the confidential computing attestation service.
func SetupCocoContainer(replicas int, image types.ImageFlags, baseImage types.ImageFlags, db shared.DbFlags) error {
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
		podman.GetServicePath(podman.ServerAttestationService+"@"), 0555, false); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf(`Environment=UYUNI_IMAGE=%s
	Environment=database_connection=jdbc:postgresql://uyuni-server.mgr.internal:%d/%s
	Environment=database_user=%s
	Environment=database_password=%s
		`, cocoImage, db.Port, db.Name, db.User, db.Password)

	if err := podman.GenerateSystemdConfFile(podman.ServerAttestationService+"@", "Service", environment); err != nil {
		return utils.Errorf(err, L("cannot generate systemd conf file"))
	}

	if err := podman.ReloadDaemon(false); err != nil {
		return err
	}

	return podman.ScaleService(replicas, podman.ServerAttestationService)
}
