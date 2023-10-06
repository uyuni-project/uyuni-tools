package shared

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/templates"
)

func GetSshAuthSocket() string {
	path := os.Getenv("SSH_AUTH_SOCK")
	if len(path) == 0 {
		log.Fatal().Msg("SSH_AUTH_SOCK is not defined, start an ssh agent and try again")
	}
	return path
}

// getSshPaths returns the user SSH config and known_hosts paths
func GetSshPaths() (string, string) {
	// Find ssh config to mount it in the container
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Msg("Failed to find home directory to look for SSH config")
	}
	sshConfigPath := filepath.Join(homedir, ".ssh", "config")
	sshKnownhostsPath := filepath.Join(homedir, ".ssh", "known_hosts")

	if !utils.FileExists(sshConfigPath) {
		sshConfigPath = ""
	}

	if !utils.FileExists(sshKnownhostsPath) {
		sshKnownhostsPath = ""
	}

	return sshConfigPath, sshKnownhostsPath
}

func GenerateMigrationScript(sourceFqdn string, kubernetes bool) string {
	scriptDir, err := os.MkdirTemp("", "uyuniadm-*")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create temporary directory")
	}

	data := templates.MigrateScriptTemplateData{
		Volumes:    utils.VOLUMES,
		SourceFqdn: sourceFqdn,
		Kubernetes: kubernetes,
	}

	scriptPath := filepath.Join(scriptDir, "migrate.sh")
	if err = utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		log.Fatal().Err(err).Msgf("Failed to generate migration script")
	}

	return scriptDir
}

func ReadTimezone(scriptDir string) string {
	data, err := os.ReadFile(filepath.Join(scriptDir, "data"))
	if err != nil {
		log.Fatal().Msgf("Failed to read data extracted from source host")
	}
	viper.SetConfigType("env")
	viper.ReadConfig(bytes.NewBuffer(data))
	return viper.GetString("Timezone")
}
