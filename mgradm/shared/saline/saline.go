// SPDX-FileCopyrightText: 2024 SUSE LLC
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
	registry string,
	salineFlags adm_utils.SalineFlags,
	baseImage types.ImageFlags,
	tz string,
	podmanArgs []string,
) error {
	if err := writeSalineServiceFiles(
		systemd, authFile, registry, salineFlags, baseImage, tz, podmanArgs,
	); err != nil {
		return err
	}

	return systemd.ScaleService(salineFlags.Replicas, podman.ServerSalineService)
}

func writeSalineServiceFiles(
	systemd podman.Systemd,
	authFile string,
	registry string,
	salineFlags adm_utils.SalineFlags,
	baseImage types.ImageFlags,
	tz string,
	podmanArgs []string,
) error {
	image := salineFlags.Image
	currentReplicas := systemd.CurrentReplicaCount(podman.ServerSalineService)
	log.Debug().Msgf("Current Saline replicas running are %d.", currentReplicas)

	if image.Tag == "" {
		if baseImage.Tag != "" {
			image.Tag = baseImage.Tag
		} else {
			image.Tag = "latest"
		}
	}
	if !salineFlags.IsChanged {
		log.Debug().Msg("Saline settings are not changed.")
		return nil
	} else if salineFlags.Replicas == 0 {
		log.Debug().Msg("No Saline requested.")
		return nil
	} else if salineFlags.Replicas > 1 {
		log.Warn().Msg(L("Multiple Saline container replicas are not currently supported, setting up only one."))
		salineFlags.Replicas = 1
	}

	salineImage, err := utils.ComputeImage(registry, baseImage.Tag, image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	pullEnabled := salineFlags.Replicas > 0 && salineFlags.IsChanged

	preparedImage, err := podman.PrepareImage(authFile, salineImage, baseImage.PullPolicy, pullEnabled)
	if err != nil {
		return err
	}

	ipv6Enabled := podman.HasIpv6Enabled(podman.UyuniNetwork)

	salineData := templates.SalineServiceTemplateData{
		NamePrefix:  "uyuni",
		Network:     podman.UyuniNetwork,
		Volumes:     utils.SalineVolumeMounts,
		Image:       preparedImage,
		SalinePort:  salineFlags.Port,
		IPV6Enabled: ipv6Enabled,
	}

	log.Info().Msg(L("Setting up Saline service"))

	if err := utils.WriteTemplateToFile(salineData,
		podman.GetServicePath(podman.ServerSalineService+"@"), 0555, true); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf(`Environment=UYUNI_SALINE_IMAGE=%s`, preparedImage)

	if err := podman.GenerateSystemdConfFile(
		podman.ServerSalineService+"@", "generated.conf", environment, true,
	); err != nil {
		return utils.Errorf(err, L("cannot generate systemd conf file"))
	}

	config := fmt.Sprintf(`Environment=TZ=%s
Environment="PODMAN_EXTRA_ARGS=%s"
`, strings.TrimSpace(tz), strings.Join(podmanArgs, " "))

	if err := podman.GenerateSystemdConfFile(podman.ServerSalineService+"@", "custom.conf",
		config, false); err != nil {
		return utils.Errorf(err, L("cannot generate systemd user configuration file"))
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
	registry string,
	salineFlags adm_utils.SalineFlags,
	baseImage types.ImageFlags,
	tz string,
	podmanArgs []string,
) error {
	if err := writeSalineServiceFiles(
		systemd, authFile, registry, salineFlags, baseImage, tz, podmanArgs,
	); err != nil {
		return err
	}
	return systemd.ScaleService(salineFlags.Replicas, podman.ServerSalineService)
}
