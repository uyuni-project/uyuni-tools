package install

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/podman"
)

func waitForSystemStart(viper *viper.Viper, globalFlags *types.GlobalFlags) {
	// Setup the systemd service configuration options
	image := fmt.Sprintf("%s:%s", viper.GetString("image"), viper.GetString("tag"))
	podman.GenerateSystemdService(viper.GetString("tz"), image, viper.GetStringSlice("podman.arg"), globalFlags.Verbose)

	log.Println("Waiting for the server to start...")
	// Start the service
	if err := exec.Command("systemctl", "enable", "--now", "uyuni-server").Run(); err != nil {
		log.Fatalf("Failed to enable uyuni-server systemd service: %s\n", err)
	}

	utils.WaitForServer()
}

func pullImage(viper *viper.Viper) {
	image := fmt.Sprintf("%s:%s", viper.GetString("image"), viper.GetString("tag"))
	log.Printf("Running podman pull %s\n", image)
	cmd := exec.Command("podman", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to pull image: %s\n", err)
	}
}

func installForPodman(viper *viper.Viper, globalFlags *types.GlobalFlags, cmd *cobra.Command, args []string) {
	pullImage(viper)

	waitForSystemStart(viper, globalFlags)

	env := map[string]string{}
	if viper.GetBool("cert.useexisting") {
		// TODO Get existing certificates path and mount them
		// Set CA_CERT, SERVER_CERT, SERVER_KEY or run the rhn-ssl-check tool in a container
		// The SERVER_CERT needs to get the intermediate keys
	} else {
		env["CERT_O"] = viper.GetString("cert.org")
		env["CERT_OU"] = viper.GetString("cert.ou")
		env["CERT_CITY"] = viper.GetString("cert.city")
		env["CERT_STATE"] = viper.GetString("cert.state")
		env["CERT_COUNTRY"] = viper.GetString("cert.country")
		env["CERT_EMAIL"] = viper.GetString("cert.email")
		env["CERT_CNAMES"] = strings.Join(append([]string{args[0]}, viper.GetStringSlice("cert.cnames")...), ",")
		env["CERT_PASS"] = viper.GetString("cert.password")
	}

	runSetup(viper, globalFlags, args[0], env)
}
