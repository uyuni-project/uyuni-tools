// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package pgsql

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// SetupPgsql prepares the systemd service and starts it if needed.
func SetupPgsql(
	systemd podman.Systemd,
	authFile string,
	pgsqlFlags cmd_utils.PgsqlFlags,
	admin string,
	password string,
) error {
	image := pgsqlFlags.Image
	currentReplicas := systemd.CurrentReplicaCount(podman.PgsqlService)
	log.Debug().Msgf("Current HUB replicas running are %d.", currentReplicas)

	if pgsqlFlags.Replicas == 0 {
		log.Debug().Msg("No pgsql requested.")
	}
	if !pgsqlFlags.IsChanged {
		log.Info().Msgf(L("No changes requested for hub. Keep %d replicas."), currentReplicas)
	}

	pullEnabled := (pgsqlFlags.Replicas > 0 && pgsqlFlags.IsChanged) || (currentReplicas > 0 && !pgsqlFlags.IsChanged)

	pgsqlImage, err := utils.ComputeImage(pgsqlFlags.Image.Registry, pgsqlFlags.Image.Tag, image)

	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	preparedImage, err := podman.PrepareImage(authFile, pgsqlImage, pgsqlFlags.Image.PullPolicy, pullEnabled)
	if err != nil {
		return err
	}

	initdbDir, _, err := utils.TempDir()
	if err != nil {
		return err
	}

	// fixme: for now, we need the script outside of this func, in EnableSSL
	// defer cleaner()

	_ = os.Chmod(initdbDir, 0555)

	data := templates.PgsqlConfigTemplateData{}

	scriptName := "pgsqlConfig.sh"
	scriptPath := filepath.Join(initdbDir, scriptName)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return fmt.Errorf(L("failed to generate %s"), scriptName)
	}

	if err := generatePgsqlSystemdService(systemd, preparedImage, initdbDir, admin, password); err != nil {
		return utils.Errorf(err, L("cannot generate systemd service"))
	}

	if err := EnablePgsql(systemd, 0); err != nil {
		return err
	}
	if err := EnablePgsql(systemd, pgsqlFlags.Replicas); err != nil {
		return err
	}
	cnx := shared.NewConnection("podman", podman.PgsqlContainerName, "")
	if err := cnx.WaitForHealthcheck(); err != nil {
		return err
	}

	// Now the servisce is up and ready, the admin credentials are no longer needed
	if err := generatePgsqlSystemdService(systemd, preparedImage, "", "", ""); err != nil {
		return utils.Errorf(err, L("cannot generate systemd service"))
	}
	return nil
}

// EnableSSL enables ssl in postgres container, as long as the certs are mounted.
func EnableSSL(systemd podman.Systemd) error {
	cnx := shared.NewConnection("podman", podman.PgsqlContainerName, "")
	if _, err := cnx.Exec("/docker-entrypoint-initdb.d/pgsqlConfig.sh"); err != nil {
		return err
	}

	if err := systemd.RestartInstantiated(podman.PgsqlService); err != nil {
		return utils.Errorf(err, L("cannot restart service"))
	}

	if err := cnx.WaitForHealthcheck(); err != nil {
		return err
	}
	return nil
}

// EnablePgsql enables the hub xmlrpc service if the number of replicas is 1.
// This function is meant for installation or migration, to enable or disable the service after, use ScaleService.
func EnablePgsql(systemd podman.Systemd, replicas int) error {
	if replicas > 1 {
		log.Warn().Msg(L("Multiple Hub XML-RPC container replicas are not currently supported, setting up only one."))
		replicas = 1
	}

	if err := systemd.ScaleService(replicas, podman.PgsqlService); err != nil {
		return utils.Errorf(err, L("cannot enable service"))
	}
	return nil
}

// Upgrade updates the systemd service files and restarts the containers if needed.
func Upgrade(
	systemd podman.Systemd,
	authFile string,
	pgsqlFlags cmd_utils.PgsqlFlags,
	admin string,
	password string,
) error {
	if err := SetupPgsql(systemd, authFile, pgsqlFlags, admin, password); err != nil {
		return err
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}

	return systemd.RestartInstantiated(podman.PgsqlService)
}

// generatePgsqlSystemdService creates the Hub XMLRPC systemd files.
func generatePgsqlSystemdService(
	systemd podman.Systemd,
	image string,
	initdb string,
	admin string,
	password string,
) error {
	pgsqlData := templates.PgsqlServiceTemplateData{
		Volumes:    utils.PgsqlRequiredVolumeMounts,
		Ports:      utils.DBPorts,
		NamePrefix: "uyuni",
		Network:    podman.UyuniNetwork,
		Image:      image,
	}
	if err := utils.WriteTemplateToFile(
		pgsqlData, podman.GetServicePath(podman.PgsqlService+"@"), 0555, true,
	); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf("Environment=UYUNI_IMAGE=%s\n", image)
	if initdb != "" {
		environment += fmt.Sprintf("Environment=PODMAN_EXTRA_ARGS=\"-v %s:/docker-entrypoint-initdb.d:z\"\n", initdb)
	}
	if admin != "" {
		environment += fmt.Sprintf("Environment=POSTGRES_USER=\"%s\"\n", admin)
	}
	if password != "" {
		environment += fmt.Sprintf("Environment=POSTGRES_PASSWORD=\"%s\"\n", password)
	}

	if err := podman.GenerateSystemdConfFile(podman.PgsqlService+"@", "generated.conf", environment, true); err != nil {
		return utils.Errorf(err, L("cannot generate systemd conf file"))
	}

	return systemd.ReloadDaemon(false)
}
