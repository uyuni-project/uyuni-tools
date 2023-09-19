package install

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/podman"
)

func waitForSystemStart(globalFlags *types.GlobalFlags, flags *InstallFlags) {
	// Setup the systemd service configuration options
	image := fmt.Sprintf("%s:%s", flags.Image.Name, flags.Image.Tag)

	podmanArgs := flags.Podman.Args
	if flags.MirrorPath != "" {
		podmanArgs = append(podmanArgs, "-v", flags.MirrorPath+":/mirror")
	}

	podman.GenerateSystemdService(flags.TZ, image, podmanArgs)

	log.Info().Msg("Waiting for the server to start...")
	// Start the service

	if err := utils.RunRawCmd("systemctl", []string{"enable", "--now", "uyuni-server"}, true); err != nil {
		log.Fatal().Err(err).Msg("Failed to enable uyuni-server systemd service")
	}

	utils.WaitForServer(globalFlags, "")
}

func pullImage(flags *InstallFlags) {
	image := fmt.Sprintf("%s:%s", flags.Image.Name, flags.Image.Tag)
	log.Info().Msgf("Running podman pull %s", image)

	err := utils.RunRawCmd("podman", []string{"pull", image}, false)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to pull image")
	}
}

func installForPodman(globalFlags *types.GlobalFlags, flags *InstallFlags, cmd *cobra.Command, args []string) {
	pullImage(flags)

	waitForSystemStart(globalFlags, flags)

	env := map[string]string{}
	if flags.Cert.UseExisting {
		// TODO Get existing certificates path and mount them
		// Set CA_CERT, SERVER_CERT, SERVER_KEY or run the rhn-ssl-check tool in a container
		// The SERVER_CERT needs to get the intermediate keys
	} else {
		env["CERT_O"] = flags.Cert.Org
		env["CERT_OU"] = flags.Cert.OU
		env["CERT_CITY"] = flags.Cert.City
		env["CERT_STATE"] = flags.Cert.State
		env["CERT_COUNTRY"] = flags.Cert.Country
		env["CERT_EMAIL"] = flags.Cert.Email
		env["CERT_CNAMES"] = strings.Join(append([]string{args[0]}, flags.Cert.Cnames...), ",")
		env["CERT_PASS"] = flags.Cert.Password
	}

	log.Info().Msg("run setup command in the container")
	runSetup(globalFlags, flags, args[0], env)

	podman.EnablePodmanSocket()
}
