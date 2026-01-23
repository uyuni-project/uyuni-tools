// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/coco"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/hub"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/pgsql"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/saline"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd shared_podman.Systemd = shared_podman.NewSystemd()

func installForPodman(
	_ *types.GlobalFlags,
	flags *podmanInstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	hostData, err := shared_podman.InspectHost()
	if err != nil {
		return err
	}

	flags.Installation.CheckParameters(cmd, "podman")
	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}

	authFile, cleaner, err := shared_podman.PodmanLogin(hostData, flags.Image.Registry, flags.Installation.SCC)
	if err != nil {
		return err
	}
	defer cleaner()

	if hostData.HasUyuniServer {
		return errors.New(
			L("Server is already initialized! Uninstall before attempting new installation or use upgrade command"),
		)
	}
	fqdn, err := utils.GetFqdn(args)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("Setting up the server with the FQDN '%s'"), fqdn)

	preparedImage, preparedPgsqlImage, err := shared_podman.PrepareImages(authFile, flags.Image, flags.Pgsql)
	if err != nil {
		return utils.Errorf(err, L("cannot prepare images"))
	}

	if err := shared_podman.SetupNetwork(false); err != nil {
		return utils.Error(err, L("cannot setup network"))
	}

	if err := podman.PrepareSSLCertificates(
		preparedImage, &flags.Installation.SSL, flags.Installation.TZ, fqdn); err != nil {
		return err
	}

	// Create all the database credentials secrets and setup the DB
	if err := setupDatabase(flags.Installation.DB, flags.Installation.ReportDB, preparedPgsqlImage); err != nil {
		return err
	}

	log.Info().Msg(L("Run setup command"))

	if err := runSetup(preparedImage, &flags.ServerFlags, fqdn); err != nil {
		return err
	}

	cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")
	if err := podman.WaitForSystemStart(systemd, cnx, preparedImage, flags.Installation.TZ,
		flags.Installation.Debug.Java, flags.Mirror, flags.Podman.Args); err != nil {
		return utils.Error(err, L("cannot wait for system start"))
	}
	if err := cnx.CopyCaCertificate(fqdn); err != nil {
		return utils.Error(err, L("failed to add SSL CA certificate to host trusted certificates"))
	}

	if path, err := exec.LookPath("uyuni-payg-extract-data"); err == nil {
		// the binary is installed
		err = utils.RunCmdStdMapping(zerolog.DebugLevel, path)
		if err != nil {
			return utils.Error(err, L("failed to extract payg data"))
		}
	}

	return utils.JoinErrors(
		shared_podman.EnablePodmanSocket(),
		coco.SetupCocoContainer(systemd, authFile, flags.Coco, flags.Image, flags.Installation.DB),
		hub.SetupHubXmlrpc(systemd, authFile, flags.Image, flags.HubXmlrpc),
		saline.SetupSalineContainer(systemd, authFile, flags.Image, flags.Saline, flags.Installation.TZ),
	)
}

func setupDatabase(dbFlags adm_utils.DBFlags, reportdbFlags adm_utils.DBFlags, preparedImage string) error {
	if err := shared_podman.CreateCredentialsSecrets(
		shared_podman.DBUserSecret, dbFlags.User,
		shared_podman.DBPassSecret, dbFlags.Password,
	); err != nil {
		return err
	}

	if err := shared_podman.CreateCredentialsSecrets(
		shared_podman.ReportDBUserSecret, reportdbFlags.User,
		shared_podman.ReportDBPassSecret, reportdbFlags.Password,
	); err != nil {
		return err
	}

	if dbFlags.IsLocal() {
		// The admin password is not needed for external databases
		if err := shared_podman.CreateCredentialsSecrets(
			shared_podman.DBAdminUserSecret, dbFlags.Admin.User,
			shared_podman.DBAdminPassSecret, dbFlags.Admin.Password,
		); err != nil {
			return err
		}

		// Run the DB container setup if the user doesn't set a custom host name for it.
		if err := pgsql.SetupPgsql(systemd, preparedImage); err != nil {
			return err
		}
	} else {
		log.Info().Msgf(
			L("Skipped database container setup to use external database %s"),
			dbFlags.Host,
		)
	}
	return nil
}

// runSetup execute the setup.
func runSetup(image string, flags *adm_utils.ServerFlags, fqdn string) error {
	env := adm_utils.GetSetupEnv(flags.Mirror, &flags.Installation, fqdn, false)
	envNames := []string{}
	envValues := []string{}
	for key, value := range env {
		envNames = append(envNames, "-e", key)
		envValues = append(envValues, fmt.Sprintf("%s=%s", key, value))
	}

	command := []string{
		"run",
		"--rm",
		"--shm-size=0",
		"--shm-size-systemd=0",
		"--name", "uyuni-setup",
		"--network", shared_podman.UyuniNetwork,
		"-e", "TZ=" + flags.Installation.TZ,
		"--secret", shared_podman.DBUserSecret + ",type=env,target=MANAGER_USER",
		"--secret", shared_podman.DBPassSecret + ",type=env,target=MANAGER_PASS",
		"--secret", shared_podman.ReportDBUserSecret + ",type=env,target=REPORT_DB_USER",
		"--secret", shared_podman.ReportDBPassSecret + ",type=env,target=REPORT_DB_PASS",
		"-e REPORT_DB_CA_CERT=" + ssl.DBCAContainerPath,
		"--secret", shared_podman.DBCASecret + ",type=mount,target=" + ssl.DBCAContainerPath,
		"--secret", shared_podman.CASecret + ",type=mount,target=" + ssl.CAContainerPath,
		"--secret", shared_podman.CASecret + ",type=mount,target=/usr/share/susemanager/salt/certs/RHN-ORG-TRUSTED-SSL-CERT",
		"--secret", shared_podman.CASecret + ",type=mount,target=/srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT",
		"--secret", shared_podman.SSLCertSecret + ",type=mount,target=" + ssl.ServerCertPath,
		"--secret", shared_podman.SSLKeySecret + ",type=mount,target=" + ssl.ServerCertKeyPath,
	}
	for _, volume := range utils.ServerVolumeMounts {
		command = append(command, "-v", fmt.Sprintf("%s:%s", volume.Name, volume.MountPath))
	}
	command = append(command, envNames...)
	command = append(command, image)

	script, err := adm_utils.GenerateSetupScript(&flags.Installation, false)
	if err != nil {
		return err
	}
	command = append(command, "/usr/bin/sh", "-e", "-c", script)

	if _, err := newRunner("podman", command...).Env(envValues).StdMapping().Exec(); err != nil {
		return utils.Error(err, L("server setup failed"))
	}

	log.Info().Msgf(L("Server set up, login on https://%[1]s with %[2]s user"), fqdn, flags.Installation.Admin.Login)
	return nil
}

var newRunner = utils.NewRunner
