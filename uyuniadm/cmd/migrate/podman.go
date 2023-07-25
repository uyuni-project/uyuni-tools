package migrate

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func migrateToPodman(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {
	sshAuthSocket := getSshAuthSocket()

	// Find ssh config to mount it in the container
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to find home directory to look for SSH config")
	}
	sshConfigPath := filepath.Join(homedir, ".ssh", "config")
	sshKnownhostsPath := filepath.Join(homedir, ".ssh", "known_hosts")

	scriptDir := generateMigrationScript(args[0], false)
	defer os.RemoveAll(scriptDir)

	extraArgs := []string{
		"-e", "SSH_AUTH_SOCK",
		"-v", filepath.Dir(sshAuthSocket) + ":" + filepath.Dir(sshAuthSocket),
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
	}

	if _, err = os.Stat(sshConfigPath); err == nil {
		extraArgs = append(extraArgs, "-v", sshConfigPath+":/root/.ssh/config")

	}

	if _, err = os.Stat(sshKnownhostsPath); err == nil {
		extraArgs = append(extraArgs, "-v", sshKnownhostsPath+":/root/.ssh/known_hosts")
	}

	log.Println("Migrating server")
	runContainer("uyuni-migration", flags.Image, flags.ImageTag, extraArgs,
		[]string{"/var/lib/uyuni-tools/migrate.sh"}, []string{}, globalFlags.Verbose)

	// Read the extracted data
	data, err := os.ReadFile(filepath.Join(scriptDir, "data"))
	if err != nil {
		log.Fatalf("Failed to read data extracted from source host")
	}
	viper.SetConfigType("env")
	viper.ReadConfig(bytes.NewBuffer(data))
	tz := viper.GetString("Timezone")

	image := fmt.Sprintf("%s:%s", flags.Image, flags.ImageTag)

	podman.GenerateSystemdService(tz, image, globalFlags.Verbose)

	// Start the service
	if err = exec.Command("systemctl", "enable", "--now", "uyuni-server").Run(); err != nil {
		log.Fatalf("Failed to enable uyuni-server systemd service: %s\n", err)
	}

	log.Println("Server migrated")
}

func runContainer(name string, image string, tag string, extraArgs []string, cmd []string, env []string, verbose bool) {

	podmanArgs := append([]string{"run"}, podman.GetCommonParams(name)...)
	podmanArgs = append(podmanArgs, extraArgs...)

	for volumeName, containerPath := range utils.VOLUMES {
		podmanArgs = append(podmanArgs, "-v", volumeName+":"+containerPath)
	}

	podmanArgs = append(podmanArgs, image+":"+tag)
	podmanArgs = append(podmanArgs, cmd...)

	podmanCmd := exec.Command("podman", podmanArgs...)

	if verbose {
		log.Printf("Running command: podman %s\n", strings.Join(podmanArgs, " "))
	}
	podmanCmd.Stdout = os.Stdout
	podmanCmd.Stderr = os.Stderr

	podmanCmd.Env = append(podmanCmd.Environ(), env...)
	if err := podmanCmd.Start(); err != nil {
		log.Fatalf("Failed to start %s container: %s\n", name, err)
	}

	// Wait for the migration to finish and report errors
	if err := podmanCmd.Wait(); err != nil {
		log.Fatalf("Failed to wait for container to finish: %s\n", err)
	}
}
