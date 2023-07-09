package migrate

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/podman"
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

	volumesArgs := []string{}
	for volumeName, containerPath := range VOLUMES {
		volumesArgs = append(volumesArgs, "-v", volumeName+":"+containerPath)
	}

	scriptDir := generateMigrationScript(args[0], false)
	defer os.RemoveAll(scriptDir)

	podmanArgs := []string{
		"run",
		"--name", "uyuni-migration",
		"--rm",
		"--tz=local",
		"--cap-add", "NET_RAW",
		"--tmpfs", "/run",
		"-v", "cgroup:/sys/fs/cgroup:rw",
		"-e", "SSH_AUTH_SOCK",
		"-v", filepath.Dir(sshAuthSocket) + ":" + filepath.Dir(sshAuthSocket),
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
	}

	if _, err = os.Stat(sshConfigPath); err == nil {
		podmanArgs = append(podmanArgs, "-v", sshConfigPath+":/root/.ssh/config")

	}

	if _, err = os.Stat(sshKnownhostsPath); err == nil {
		podmanArgs = append(podmanArgs, "-v", sshKnownhostsPath+":/root/.ssh/known_hosts")
	}

	podmanArgs = append(podmanArgs, volumesArgs...)

	podmanArgs = append(podmanArgs,
		flags.Image+":"+flags.ImageTag,
		"/var/lib/uyuni-tools/migrate.sh",
	)

	log.Println("Migrating server")

	podmanCmd := exec.Command("podman", podmanArgs...)

	if globalFlags.Verbose {
		log.Printf("Running command: podman %s\n", strings.Join(podmanArgs, " "))
		podmanCmd.Stdout = os.Stdout
		podmanCmd.Stderr = os.Stderr
	}

	if err = podmanCmd.Start(); err != nil {
		log.Fatalf("Failed to start migration container: %s\n", err)
	}

	// Wait for the migration to finish and report errors
	if err = podmanCmd.Wait(); err != nil {
		log.Fatalf("Failed to wait for container to finish: %s\n", err)
	}

	// Setup the systemd service configuration options
	config := podman.ReadConfig()

	config = podman.UpdateConfigValue(config, "NAMESPACE", filepath.Dir(flags.Image))
	config = podman.UpdateConfigValue(config, "TAG", flags.ImageTag)

	// TODO More values to write like UYUNI_FQDN?
	podman.WriteConfig(config)

	if globalFlags.Verbose {
		log.Printf("Wrote configuration:\n%s\n", config)
	}

	// Start the service
	if err = exec.Command("systemctl", "enable", "--now", "uyuni-server").Run(); err != nil {
		log.Fatalf("Failed to enable uyuni-server systemd service: %s\n", err)
	}

	log.Println("Server migrated")
}
