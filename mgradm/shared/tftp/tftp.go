// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package tftp

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

// SetupTFTPContainer() prepares the systemd service for the TFTP server and starts it if needed.
// tag is the global images tag.
func SetupTFTPContainer(
	systemd podman.Systemd,
	authFile string,
	baseImage types.ImageFlags,
	tftpFlags adm_utils.TFTPDFlags,
	fqdn string,
) error {
	if tftpFlags.Disable && systemd.ServiceIsEnabled(podman.TFTPService) {
		log.Debug().Msgf("The TFTP service is no longer requested")
		if err := systemd.DisableService(podman.TFTPService); err != nil {
			return err
		}
	}

	tftpImage, err := utils.ComputeImage(baseImage.Registry.Host, baseImage.Tag, tftpFlags.Image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	preparedImage, err := podman.PrepareImage(authFile, tftpImage, baseImage.PullPolicy, !tftpFlags.Disable)
	if err != nil {
		return err
	}

	if err := generateTFTPSystemdService(systemd, preparedImage, fqdn); err != nil {
		return utils.Errorf(err, L("cannot generate systemd service"))
	}

	if !tftpFlags.Disable {
		if err := systemd.EnableService(podman.TFTPService); err != nil {
			return err
		}
	}
	return nil
}

// Upgrade updates the systemd service files and restarts the container if needed.
func Upgrade(
	systemd podman.Systemd,
	authFile string,
	baseImage types.ImageFlags,
	tftpFlags adm_utils.TFTPDFlags,
	fqdn string,
) error {
	if tftpFlags.Image.Name == "" {
		// Don't touch the tftp service in ptf if not already present.
		return nil
	}
	if err := SetupTFTPContainer(systemd, authFile, baseImage, tftpFlags, fqdn); err != nil {
		return err
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}

	if !tftpFlags.IsChanged {
		return systemd.RestartInstantiated(podman.HubXmlrpcService)
	}
	if systemd.ServiceIsEnabled(podman.TFTPService) {
		return systemd.RestartService(podman.TFTPService)
	}
	return systemd.EnableService(podman.TFTPService)
}

// generateTFTPSystemdService creates the TFTP systemd files.
func generateTFTPSystemdService(systemd podman.Systemd, image string, fqdn string) error {
	tftpData := templates.TFTPDTemplateData{
		CaSecret:   podman.CASecret,
		Network:    podman.UyuniNetwork,
		ServerFQDN: fqdn,
	}
	if err := utils.WriteTemplateToFile(
		tftpData, podman.GetServicePath(podman.TFTPService), 0555, true,
	); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf("Environment=UYUNI_TFTPD_IMAGE=%s", image)
	if err := podman.GenerateSystemdConfFile(
		podman.TFTPService, "generated.conf", environment, true,
	); err != nil {
		return utils.Errorf(err, L("cannot generate systemd conf file"))
	}

	return systemd.ReloadDaemon(false)
}
