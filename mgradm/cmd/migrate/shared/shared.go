// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
)

func GetSshAuthSocket() string {
	path := os.Getenv("SSH_AUTH_SOCK")
	if len(path) == 0 {
		log.Fatal().Msg("SSH_AUTH_SOCK is not defined, start an ssh agent and try again")
	}
	return path
}

// GetSshPaths returns the user SSH config and known_hosts paths
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

// GetCustomSELinuxPolicyDetails returns the custom SELinux policy path and the Podman label
func GetCustomSELinuxPolicyDetails(productName string) (string, string) {
	podmanLabel := "label=type:" + productName + "-selinux-policy.process"

	fileName := productName + "-selinux-policy.cil"
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Msgf("Failed to find home directory to look for %s", fileName)
	}
	filePath := filepath.Join(homedir, fileName)

	return podmanLabel, filePath
}

// InstallCustomSELinuxPolicy make use of semodule command to install the custom policy
func InstallCustomSELinuxPolicy(policyPath string) {
	if !utils.FileExists(policyPath) {
		log.Fatal().Msgf("Failed to load the SELinux policy: %s", policyPath)
	}

	udicaTemplatesPath := "/usr/share/udica/templates"
	if !utils.FileExists(udicaTemplatesPath) {
		log.Fatal().Msgf("Udica SELinux templates are not present on: %s\n"+
			"Please install Udica before continue: https://github.com/containers/udica",
			udicaTemplatesPath)
	}

	udicaBaseContainerPolicy := filepath.Join(udicaTemplatesPath, "base_container.cil")
	udicaHomeContainerPolicy := filepath.Join(udicaTemplatesPath, "home_container.cil")
	udicaTmpContainerPolicy := filepath.Join(udicaTemplatesPath, "tmp_container.cil")
	udicaNetContainerPolicy := filepath.Join(udicaTemplatesPath, "net_container.cil")
	udicaLogContainerPolicy := filepath.Join(udicaTemplatesPath, "log_container.cil")
	udicaConfigContainerPolicy := filepath.Join(udicaTemplatesPath, "config_container.cil")
	udicaTTYContainerPolicy := filepath.Join(udicaTemplatesPath, "tty_container.cil")
	udicaVirtContainerPolicy := filepath.Join(udicaTemplatesPath, "virt_container.cil")
	udicaXContainerPolicy := filepath.Join(udicaTemplatesPath, "x_container.cil")

	errInstall := utils.RunCmdStdMapping("semodule",
		"-i", policyPath, udicaBaseContainerPolicy,
		udicaHomeContainerPolicy, udicaTmpContainerPolicy,
		udicaNetContainerPolicy, udicaLogContainerPolicy,
		udicaConfigContainerPolicy, udicaTTYContainerPolicy,
		udicaVirtContainerPolicy, udicaXContainerPolicy)

	if errInstall != nil {
		log.Fatal().Err(errInstall).Msg("Custom SELinux policies can't be installed.")
	}
}

func GenerateMigrationScript(sourceFqdn string, kubernetes bool) string {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
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

func GeneratePgMigrationScript(scriptDir string, oldPgVersion string, newPgVersion string, kubernetes bool) {
	data := templates.MigratePostgresVersionTemplateData {
		OldVersion: oldPgVersion,
		NewVersion: newPgVersion,
		Kubernetes: kubernetes,
	}

	scriptPath := filepath.Join(scriptDir, "migrate.sh")
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		log.Fatal().Err(err).Msgf("Failed to generate migration script")
	}
}

func GenerateFinalizePostgresMigrationScript(scriptDir string, RunAutotune bool, RunReindex bool, RunSchemaUpdate bool, RunDistroMigration bool, kubernetes bool) {
	data := templates.FinalizePostgresTemplateData {
		RunAutotune: RunAutotune,
		RunReindex: RunReindex,
		RunSchemaUpdate: RunSchemaUpdate,
		RunDistroMigration: RunDistroMigration,
		Kubernetes: kubernetes,
	}

	scriptPath := filepath.Join(scriptDir, "migrate.sh")
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		log.Fatal().Err(err).Msgf("Failed to generate migration script")
	}
}

func ReadContainerData(scriptDir string) (string, string, string) {
	data, err := os.ReadFile(filepath.Join(scriptDir, "data"))
	if err != nil {
		log.Fatal().Msgf("Failed to read data extracted from source host")
	}
	viper.SetConfigType("env")
	viper.ReadConfig(bytes.NewBuffer(data))
	return viper.GetString("Timezone"), viper.GetString("old_pg_version"), viper.GetString("new_pg_version")
}
