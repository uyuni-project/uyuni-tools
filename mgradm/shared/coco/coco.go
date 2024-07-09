// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package coco

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Upgrade coco attestation.
func Upgrade(image types.ImageFlags, baseImage types.ImageFlags, dbPort int, dbName string, dbUser string, dbPassword string) error {
	if err := podman.StopInstantiated(podman.ServerAttestationService); err != nil {
		return err
	}
	if err := writeCocoServiceFiles(image, baseImage, dbName, dbPort, dbUser, dbPassword); err != nil {
		return err
	}
	return podman.StartInstantiated(podman.ServerAttestationService)
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

	inspectedHostValues, err := utils.InspectHost(false)
	if err != nil {
		return utils.Errorf(err, L("cannot inspect host values"))
	}

	pullArgs := []string{}
	_, scc_user_exist := inspectedHostValues["host_scc_username"]
	_, scc_user_password := inspectedHostValues["host_scc_password"]
	if scc_user_exist && scc_user_password {
		pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
	}

	preparedImage, err := podman.PrepareImage(cocoImage, baseImage.PullPolicy, pullArgs...)
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
	Environment=database_user=%s
	Environment=database_password=%s`, preparedImage, dbPort, dbName, dbUser, dbPassword)

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
