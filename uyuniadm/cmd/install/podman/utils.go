package podman

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/podman"
)

func waitForSystemStart(cnx *utils.Connection, flags *podmanInstallFlags) {
	// Setup the systemd service configuration options
	image := fmt.Sprintf("%s:%s", flags.Image.Name, flags.Image.Tag)

	podmanArgs := flags.Podman.Args
	if flags.MirrorPath != "" {
		podmanArgs = append(podmanArgs, "-v", flags.MirrorPath+":/mirror")
	}

	podman.GenerateSystemdService(flags.TZ, image, flags.Debug.Java, podmanArgs)

	log.Info().Msg("Waiting for the server to start...")
	// Start the service

	if err := utils.RunCmd("systemctl", "enable", "--now", "uyuni-server"); err != nil {
		log.Fatal().Err(err).Msg("Failed to enable uyuni-server systemd service")
	}

	cnx.WaitForServer()
}

func pullImage(flags *podmanInstallFlags) {
	image := fmt.Sprintf("%s:%s", flags.Image.Name, flags.Image.Tag)
	log.Info().Msgf("Running podman pull %s", image)

	err := utils.RunCmdStdMapping("podman", "pull", image)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to pull image")
	}
}

func installForPodman(globalFlags *types.GlobalFlags, flags *podmanInstallFlags, cmd *cobra.Command, args []string) {
	fqdn := getFqdn(args)
	log.Info().Msgf("setting up server with the FQDN '%s'", fqdn)

	pullImage(flags)

	cnx := utils.NewConnection("podman")
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

	shared.RunSetup(cnx, &flags.InstallFlags, fqdn, env)

	if flags.Ssl.UseExisting() {
		podman.UpdateSslCertificate(cnx, &flags.Ssl.Ca, &flags.Ssl.Server)
	}

	podman.EnablePodmanSocket()
}

func getFqdn(args []string) string {
	if len(args) == 1 {
		return args[0]
	} else {
		fqdn_b, err := utils.RunCmdOutput(zerolog.DebugLevel, "hostname", "-f")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to compute server FQDN")
		}
		return string(fqdn_b)
	}
}
