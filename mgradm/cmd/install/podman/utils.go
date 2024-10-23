// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	install_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/coco"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/hub"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func waitForSystemStart(cnx *shared.Connection, image string, flags *podmanInstallFlags) error {
	err := podman.GenerateSystemdService(flags.TZ, image, flags.Debug.Java, flags.Mirror, flags.Podman.Args)
	if err != nil {
		return err
	}

	log.Info().Msg(L("Waiting for the server to start…"))
	if err := shared_podman.EnableService(shared_podman.ServerService); err != nil {
		return utils.Errorf(err, L("cannot enable service"))
	}

	return cnx.WaitForServer()
}

func installForPodman(
	globalFlags *types.GlobalFlags,
	flags *podmanInstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	hostData, err := shared_podman.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := shared_podman.PodmanLogin(hostData, flags.Scc)
	if err != nil {
		return utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	if hostData.HasUyuniServer {
		return errors.New(
			L("Server is already initialized! Uninstall before attempting new installation or use upgrade command"),
		)
	}

	flags.CheckParameters(cmd, "podman")
	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}

	fqdn, err := utils.GetFqdn(args)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("Setting up the server with the FQDN '%s'"), fqdn)

	image, err := utils.ComputeImage(flags.Image.Registry, utils.DefaultTag, flags.Image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	preparedImage, err := shared_podman.PrepareImage(authFile, image, flags.Image.PullPolicy, true)
	if err != nil {
		return err
	}

	cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")
	if err := waitForSystemStart(cnx, preparedImage, flags); err != nil {
		return utils.Errorf(err, L("cannot wait for system start"))
	}

	caPassword := flags.Ssl.Password
	if flags.Ssl.UseExisting() {
		// We need to have a password for the generated CA, even though it will be thrown away after install
		caPassword = "dummy"
	}

	env := map[string]string{
		"CERT_O":       flags.Ssl.Org,
		"CERT_OU":      flags.Ssl.OU,
		"CERT_CITY":    flags.Ssl.City,
		"CERT_STATE":   flags.Ssl.State,
		"CERT_COUNTRY": flags.Ssl.Country,
		"CERT_EMAIL":   flags.Ssl.Email,
		"CERT_CNAMES":  strings.Join(append([]string{fqdn}, flags.Ssl.Cnames...), ","),
		"CERT_PASS":    caPassword,
	}

	log.Info().Msg(L("Run setup command in the container"))

	if err := install_shared.RunSetup(cnx, &flags.InstallFlags, fqdn, env); err != nil {
		if stopErr := shared_podman.StopService(shared_podman.ServerService); stopErr != nil {
			log.Error().Msgf(L("Failed to stop service: %v"), stopErr)
		}
		return err
	}

	if path, err := exec.LookPath("uyuni-payg-extract-data"); err == nil {
		// the binary is installed
		err = utils.RunCmdStdMapping(zerolog.DebugLevel, path)
		if err != nil {
			return utils.Errorf(err, L("failed to extract payg data"))
		}
	}

	if flags.Coco.Replicas > 0 {
		// This may need to be moved up later once more containers require DB access
		if err := shared_podman.CreateDbSecrets(flags.Db.User, flags.Db.Password); err != nil {
			return err
		}
		if err := coco.SetupCocoContainer(
			authFile, flags.Image.Registry, flags.Coco, flags.Image,
			flags.Db.Name, flags.Db.Port,
		); err != nil {
			return err
		}
	}

	if flags.HubXmlrpc.Replicas > 0 {
		if err := hub.SetupHubXmlrpc(
			authFile, flags.Image.Registry, flags.Image.PullPolicy, flags.Image.Tag, flags.HubXmlrpc,
		); err != nil {
			return err
		}
	}

	if flags.Ssl.UseExisting() {
		if err := podman.UpdateSslCertificate(cnx, &flags.Ssl.Ca, &flags.Ssl.Server); err != nil {
			return utils.Errorf(err, L("cannot update SSL certificate"))
		}
	}

	if err := shared_podman.EnablePodmanSocket(); err != nil {
		return utils.Errorf(err, L("cannot enable podman socket"))
	}
	return nil
}
