// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	install_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func waitForSystemStart(cnx *shared.Connection, flags *podmanInstallFlags) {
	// Setup the systemd service configuration options
	image := fmt.Sprintf("%s:%s", flags.Image.Name, flags.Image.Tag)

	podmanArgs := flags.Podman.Args
	if flags.MirrorPath != "" {
		podmanArgs = append(podmanArgs, "-v", flags.MirrorPath+":/mirror")
	}

	podman.GenerateSystemdService(flags.TZ, image, flags.Debug.Java, podmanArgs)

	log.Info().Msg("Waiting for the server to start...")
	shared_podman.EnableService("uyuni-server")

	cnx.WaitForServer()
}

func installForPodman(
	globalFlags *types.GlobalFlags,
	flags *podmanInstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	flags.CheckParameters(cmd, "podman")
	if _, err := exec.LookPath("podman"); err != nil {
		return fmt.Errorf("install podman before running this command: %s", err)
	}

	fqdn, err := getFqdn(args)
	if err != nil {
		return err
	}
	log.Info().Msgf("setting up server with the FQDN '%s'", fqdn)

	image, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		return fmt.Errorf("failed to compute image URL, %s", err)
	}
	shared_podman.PrepareImage(image, flags.Image.PullPolicy)

	cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")
	waitForSystemStart(cnx, flags)

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

	log.Info().Msg("run setup command in the container")

	if err := install_shared.RunSetup(cnx, &flags.InstallFlags, fqdn, env); err != nil {
		return err
	}

	if flags.Ssl.UseExisting() {
		podman.UpdateSslCertificate(cnx, &flags.Ssl.Ca, &flags.Ssl.Server)
	}

	shared_podman.EnablePodmanSocket()
	return nil
}

func getFqdn(args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	} else {
		fqdn_b, err := utils.RunCmdOutput(zerolog.DebugLevel, "hostname", "-f")
		if err != nil {
			return "", fmt.Errorf("failed to compute server FQDN: %s", err)
		}
		return strings.TrimSpace(string(fqdn_b)), nil
	}
}
