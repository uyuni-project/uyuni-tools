package migrate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/podman"
)

func migrateToPodman(globalFlags *types.GlobalFlags, flags *MigrateFlags, cmd *cobra.Command, args []string) {
	// Find the SSH Socket and paths for the migration
	sshAuthSocket := getSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := getSshPaths()

	scriptDir := generateMigrationScript(args[0], false)
	defer os.RemoveAll(scriptDir)

	extraArgs := []string{
		"-e", "SSH_AUTH_SOCK",
		"-v", filepath.Dir(sshAuthSocket) + ":" + filepath.Dir(sshAuthSocket),
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
	}

	if sshConfigPath != "" {
		extraArgs = append(extraArgs, "-v", sshConfigPath+":/root/.ssh/config")
	}

	if sshKnownhostsPath != "" {
		extraArgs = append(extraArgs, "-v", sshKnownhostsPath+":/root/.ssh/known_hosts")
	}

	log.Info().Msg("Migrating server")
	runContainer("uyuni-migration", flags.Image.Name, flags.Image.Tag, extraArgs,
		[]string{"/var/lib/uyuni-tools/migrate.sh"})

	// Read the extracted data
	tz := readTimezone(scriptDir)
	fullImage := fmt.Sprintf("%s:%s", flags.Image.Name, flags.Image.Tag)

	podman.GenerateSystemdService(tz, fullImage, viper.GetStringSlice("podman.arg"))

	// Start the service

	if err := utils.RunRawCmd("systemctl", []string{"enable", "--now", "uyuni-server"}, true); err != nil {
		log.Fatal().Err(err).Msgf("Failed to enable uyuni-server systemd service")
	}

	log.Info().Msg("Server migrated")

	podman.EnablePodmanSocket()
}

func runContainer(name string, image string, tag string, extraArgs []string, cmd []string) {

	podmanArgs := append([]string{"run"}, podman.GetCommonParams(name)...)
	podmanArgs = append(podmanArgs, extraArgs...)

	for volumeName, containerPath := range utils.VOLUMES {
		podmanArgs = append(podmanArgs, "-v", volumeName+":"+containerPath)
	}

	podmanArgs = append(podmanArgs, image+":"+tag)
	podmanArgs = append(podmanArgs, cmd...)

	err := utils.RunRawCmd("podman", podmanArgs, false)

	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to run %s container", name)
	}
}
