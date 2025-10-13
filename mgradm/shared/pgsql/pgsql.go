// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package pgsql

import (
	"fmt"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func PreparePgsqlImage(
	authFile string,
	pgsqlFlags *types.PgsqlFlags,
	globalImageFlags *types.ImageFlags,
) (string, error) {
	image := pgsqlFlags.Image
	pgsqlImage, err := utils.ComputeImage(globalImageFlags.Registry.Host, globalImageFlags.Tag, image)

	if err != nil {
		return "", utils.Error(err, L("failed to compute image URL"))
	}

	preparedImage, err := podman.PrepareImage(authFile, pgsqlImage, globalImageFlags.PullPolicy, true)
	if err != nil {
		return "", err
	}
	return preparedImage, err
}

// SetupPgsql prepares the systemd service and starts it if needed.
func SetupPgsql(
	systemd podman.Systemd,
	pgsqlImage string,
) error {
	if err := GeneratePgsqlSystemdService(systemd, pgsqlImage); err != nil {
		return utils.Error(err, L("cannot generate systemd service"))
	}

	if err := systemd.EnableService(podman.DBService); err != nil {
		return err
	}
	if err := systemd.StartService(podman.DBService); err != nil {
		return err
	}
	cnx := shared.NewConnection("podman", podman.DBContainerName, "")
	if err := cnx.WaitForHealthcheck(); err != nil {
		return utils.Errorf(err, L("%s fails healtcheck"), podman.DBContainerName)
	}

	return nil
}

// Upgrade updates the systemd service files and restarts the containers if needed.
func Upgrade(
	preparedImage string,
	systemd podman.Systemd,
) error {
	if err := GeneratePgsqlSystemdService(systemd, preparedImage); err != nil {
		return utils.Error(err, L("cannot generate systemd service"))
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}

	if err := systemd.EnableService(podman.DBService); err != nil {
		return err
	}

	if err := systemd.StartService(podman.DBService); err != nil {
		return err
	}

	cnx := shared.NewConnection("podman", podman.DBContainerName, "")
	if err := cnx.WaitForHealthcheck(); err != nil {
		return utils.Errorf(err, L("%s fails healtcheck"), podman.DBContainerName)
	}

	return nil
}

// GeneratePgsqlSystemdService creates the DB container systemd files.
func GeneratePgsqlSystemdService(
	systemd podman.Systemd,
	image string,
) error {
	pgsqlData := templates.PgsqlServiceTemplateData{
		Volumes:         utils.PgsqlRequiredVolumeMounts,
		Ports:           utils.DBPorts,
		NamePrefix:      "uyuni",
		Network:         podman.UyuniNetwork,
		Image:           image,
		CaSecret:        podman.DBCASecret,
		CaPath:          ssl.DBCAContainerPath,
		CertSecret:      podman.DBSSLCertSecret,
		CertPath:        ssl.DBCertPath,
		KeySecret:       podman.DBSSLKeySecret,
		KeyPath:         ssl.DBCertKeyPath,
		AdminUser:       podman.DBAdminUserSecret,
		AdminPassword:   podman.DBAdminPassSecret,
		ManagerUser:     podman.DBUserSecret,
		ManagerPassword: podman.DBPassSecret,
		ReportUser:      podman.ReportDBUserSecret,
		ReportPassword:  podman.ReportDBPassSecret,
	}
	if err := utils.WriteTemplateToFile(
		pgsqlData, podman.GetServicePath(podman.DBService), 0555, true,
	); err != nil {
		return utils.Error(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf("Environment=UYUNI_IMAGE=%s\n", image)

	if err := podman.GenerateSystemdConfFile(podman.DBService, "generated.conf", environment, true); err != nil {
		return utils.Error(err, L("cannot generate systemd configuration file"))
	}

	return systemd.ReloadDaemon(false)
}
