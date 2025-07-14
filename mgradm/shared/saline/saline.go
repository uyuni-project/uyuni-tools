// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package saline

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Upgrade Saline.
func Upgrade(
	systemd podman.Systemd,
	authFile string,
	salineFlags adm_utils.SalineFlags,
	baseImage types.ImageFlags,
	tz string,
) error {
	if err := writeSalineServiceFiles(
		systemd, authFile, salineFlags, baseImage, tz,
	); err != nil {
		return err
	}

	if salineFlags.Replicas > 0 {
		if systemd.IsServiceRunning(podman.SalineService) {
			return systemd.RestartService(podman.SalineService)
		}
		return systemd.EnableService(podman.SalineService)
	}
	if systemd.ServiceIsEnabled(podman.SalineService) {
		return systemd.DisableService(podman.SalineService)
	}
	return nil
}

func writeSalineServiceFiles(
	systemd podman.Systemd,
	authFile string,
	salineFlags adm_utils.SalineFlags,
	baseImage types.ImageFlags,
	tz string,
) error {
	image := salineFlags.Image

	if image.Tag == "" {
		if baseImage.Tag != "" {
			image.Tag = baseImage.Tag
		} else {
			image.Tag = "latest"
		}
	}
	if !salineFlags.IsChanged {
		log.Debug().Msg("Saline settings are not changed.")
	} else if salineFlags.Replicas == 0 {
		log.Debug().Msg("No Saline requested.")
	} else if salineFlags.Replicas > 1 {
		log.Warn().Msg(L("Multiple Saline container replicas are not currently supported, setting up only one."))
		salineFlags.Replicas = 1
	}

	salineImage, err := utils.ComputeImage(image)
	if err != nil {
		return utils.Error(err, L("failed to compute image URL"))
	}

	pullEnabled := salineFlags.Replicas > 0 && salineFlags.IsChanged

	preparedImage, err := podman.PrepareImage(authFile, salineImage, baseImage.PullPolicy, pullEnabled)
	if err != nil {
		return err
	}

	salineData := templates.SalineServiceTemplateData{
		NamePrefix: "uyuni",
		Network:    podman.UyuniNetwork,
		Volumes:    utils.SalineVolumeMounts,
		Image:      preparedImage.Name,
	}

	log.Info().Msg(L("Setting up Saline service"))

	if err := utils.WriteTemplateToFile(salineData,
		podman.GetServicePath(podman.SalineService+"@"), 0555, true); err != nil {
		return utils.Error(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf(`Environment=UYUNI_SALINE_IMAGE=%s`, preparedImage.Name)

	if err := podman.GenerateSystemdConfFile(
		podman.SalineService+"@", "generated.conf", environment, true,
	); err != nil {
		return utils.Error(err, L("cannot generate systemd conf file"))
	}

	config := fmt.Sprintf(`Environment=TZ=%s
`, strings.TrimSpace(tz))

	if err := podman.GenerateSystemdConfFile(podman.SalineService+"@", "custom.conf",
		config, false); err != nil {
		return utils.Error(err, L("cannot generate systemd user configuration file"))
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}
	return nil
}

// SetupSalineContainer sets up the Saline service.
func SetupSalineContainer(
	systemd podman.Systemd,
	authFile string,
	salineFlags adm_utils.SalineFlags,
	baseImage types.ImageFlags,
	tz string,
) error {
	if err := writeSalineServiceFiles(systemd, authFile, salineFlags, baseImage, tz); err != nil {
		return err
	}
	return EnableSaline(systemd, salineFlags.Replicas)
}

// EnableSaline enables the saline service if the number of replicas is 1.
// This function is meant for installation or migration, to enable or disable the service after, use ScaleService.
func EnableSaline(systemd podman.Systemd, replicas int) error {
	if replicas > 1 {
		log.Warn().Msg(L("Multiple Saline container replicas are not currently supported, setting up only one."))
		replicas = 1
	}

	if replicas > 0 {
		if err := systemd.ScaleService(replicas, podman.SalineService); err != nil {
			return utils.Errorf(err, L("cannot enable service"))
		}
	}
	return nil
}
