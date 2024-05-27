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
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func setupHubXmlrpcContainer(flags *podmanInstallFlags) error {
	if flags.HubXmlrpc.Enable {
		log.Info().Msg(L("Enabling Hub XML-RPC API container."))
		tag := flags.HubXmlrpc.Image.Tag
		if tag == "" {
			tag = flags.Image.Tag
		}
		hubXmlrpcImage, err := utils.ComputeImage(flags.HubXmlrpc.Image.Name, tag)
		if err != nil {
			return utils.Errorf(err, L("failed to compute image URL"))
		}

		if err := podman.GenerateHubXmlrpcSystemdService(hubXmlrpcImage); err != nil {
			return utils.Errorf(err, L("cannot generate systemd service"))
		}

		if err := shared_podman.EnableService(shared_podman.HubXmlrpcService); err != nil {
			return utils.Errorf(err, L("cannot enable service"))
		}
	}
	return nil
}

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
	flags.CheckParameters(cmd, "podman")
	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}

	inspectedHostValues, err := utils.InspectHost(false)
	if err != nil {
		return utils.Errorf(err, L("cannot inspect host values"))
	}

	fqdn, err := getFqdn(args)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("Setting up the server with the FQDN '%s'"), fqdn)

	image, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}
	pullArgs := []string{}
	_, scc_user_exist := inspectedHostValues["host_scc_username"]
	_, scc_user_password := inspectedHostValues["host_scc_password"]
	if scc_user_exist && scc_user_password {
		pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
	}

	preparedImage, err := shared_podman.PrepareImage(image, flags.Image.PullPolicy, pullArgs...)
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

	if err := coco.SetupCocoContainer(flags.Coco.Replicas, flags.Coco.Image, flags.Image, flags.Db); err != nil {
		return err
	}

	if err := setupHubXmlrpcContainer(flags); err != nil {
		return err
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

func getFqdn(args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	} else {
		fqdn_b, err := utils.RunCmdOutput(zerolog.DebugLevel, "hostname", "-f")
		if err != nil {
			return "", utils.Errorf(err, L("failed to compute server FQDN"))
		}
		return strings.TrimSpace(string(fqdn_b)), nil
	}
}
