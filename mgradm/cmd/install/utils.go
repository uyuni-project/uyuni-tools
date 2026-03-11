// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"errors"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/coco"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/hub"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/pgsql"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/saline"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/tftp"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
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

	if err := podman.GenerateSystemdService(
		systemd, preparedImage, flags.Installation, flags.Podman.Args, flags.Mirror, fqdn,
	); err != nil {
		return utils.Error(err, L("failed to generate server service"))
	}

	cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")
	if err := podman.WaitForSystemStart(systemd, cnx, preparedImage, flags.Installation.TZ,
		flags.Installation.Debug.Java, flags.Mirror, flags.Podman.Args, fqdn); err != nil {
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
		tftp.SetupTFTPContainer(systemd, authFile, flags.Image, flags.TFTPD, fqdn),
	)
}

func setupDatabase(dbFlags types.DBFlags, reportdbFlags types.DBFlags, preparedImage string) error {
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
