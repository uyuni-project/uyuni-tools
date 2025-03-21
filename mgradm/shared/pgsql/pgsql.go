// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package pgsql

import (
	"fmt"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// SetupPgsql prepares the systemd service and starts it if needed.
func SetupPgsql(
	systemd podman.Systemd,
	authFile string,
	pgsqlFlags *cmd_utils.PgsqlFlags,
	globalImageFlags *types.ImageFlags,
) error {
	image := pgsqlFlags.Image
	pgsqlImage, err := utils.ComputeImage(globalImageFlags.Registry, globalImageFlags.Tag, image)

	if err != nil {
		return utils.Error(err, L("failed to compute image URL"))
	}

	preparedImage, err := podman.PrepareImage(authFile, pgsqlImage, globalImageFlags.PullPolicy, true)
	if err != nil {
		return err
	}

	if err := generatePgsqlSystemdService(systemd, preparedImage); err != nil {
		return utils.Error(err, L("cannot generate systemd service"))
	}

	if err := EnablePgsql(systemd); err != nil {
		return err
	}
	cnx := shared.NewConnection("podman", podman.DBContainerName, "")
	if err := cnx.WaitForHealthcheck(); err != nil {
		return err
	}

	return nil
}

// EnablePgsql enables the database service.
// This function is meant for installation or migration, to enable and start the service.
func EnablePgsql(systemd podman.Systemd) error {
	if err := systemd.EnableService(podman.DBService); err != nil {
		return utils.Errorf(err, L("cannot enable %s service"), podman.DBService)
	}
	return nil
}

// Upgrade updates the systemd service files and restarts the containers if needed.
func Upgrade(
	systemd podman.Systemd,
	authFile string,
	pgsqlFlags cmd_utils.PgsqlFlags,
) error {
	image := pgsqlFlags.Image
	pgsqlImage, err := utils.ComputeImage(pgsqlFlags.Image.Registry, pgsqlFlags.Image.Tag, image)

	if err != nil {
		return utils.Error(err, L("failed to compute image URL"))
	}

	preparedImage, err := podman.PrepareImage(authFile, pgsqlImage, pgsqlFlags.Image.PullPolicy, true)
	if err != nil {
		return err
	}

	if err := generatePgsqlSystemdService(systemd, preparedImage); err != nil {
		return utils.Error(err, L("cannot generate systemd service"))
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}

	if err := EnablePgsql(systemd); err != nil {
		return err
	}

	cnx := shared.NewConnection("podman", podman.DBContainerName, "")
	return cnx.WaitForHealthcheck()
}

// generatePgsqlSystemdService creates the DB container systemd files.
func generatePgsqlSystemdService(
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
		CertSecret:      podman.DBSSLCertSecret,
		KeySecret:       podman.DBSSLKeySecret,
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
