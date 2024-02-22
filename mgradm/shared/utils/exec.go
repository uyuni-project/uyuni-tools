// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ExecCommand execute commands passed as argument in the current system.
func ExecCommand(logLevel zerolog.Level, cnx *shared.Connection, args ...string) error {
	podName, err := cnx.GetPodName()
	if err != nil {
		return fmt.Errorf("execCommand failed %s", err)
	}

	commandArgs := []string{"exec", podName}

	command, err := cnx.GetCommand()
	if err != nil {
		log.Fatal().Err(err)
	}

	if command == "kubectl" {
		commandArgs = append(commandArgs, "-c", "uyuni", "--")
	}

	commandArgs = append(commandArgs, "sh", "-c", strings.Join(args, " "))

	runCmd := exec.Command(command, commandArgs...)
	logger := utils.OutputLogWriter{Logger: log.Logger, LogLevel: logLevel}
	runCmd.Stdout = logger
	runCmd.Stderr = logger
	return runCmd.Run()
}

// GeneratePgsqlVersionUpgradeScript generates the PostgreSQL version upgrade script.
func GeneratePgsqlVersionUpgradeScript(scriptDir string, oldPgVersion string, newPgVersion string, kubernetes bool) (string, error) {
	data := templates.PostgreSQLVersionUpgradeTemplateData{
		OldVersion: oldPgVersion,
		NewVersion: newPgVersion,
		Kubernetes: kubernetes,
	}

	scriptName := "pgsqlVersionUpgrade.sh"
	scriptPath := filepath.Join(scriptDir, scriptName)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return "", fmt.Errorf("failed to generate %s", scriptName)
	}
	return scriptName, nil
}

// GenerateFinalizePostgresScript generates the script to finalize PostgreSQL upgrade.
func GenerateFinalizePostgresScript(scriptDir string, RunAutotune bool, RunReindex bool, RunSchemaUpdate bool, RunDistroMigration bool, kubernetes bool) (string, error) {
	data := templates.FinalizePostgresTemplateData{
		RunAutotune:        RunAutotune,
		RunReindex:         RunReindex,
		RunSchemaUpdate:    RunSchemaUpdate,
		RunDistroMigration: RunDistroMigration,
		Kubernetes:         kubernetes,
	}

	scriptName := "pgsqlFinalize.sh"
	scriptPath := filepath.Join(scriptDir, scriptName)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return "", fmt.Errorf("failed to generate %s", scriptName)
	}
	return scriptName, nil
}

// GeneratePostUpgradeScript generates the script to be run after upgrade.
func GeneratePostUpgradeScript(scriptDir string, cobblerHost string) (string, error) {
	data := templates.PostUpgradeTemplateData{
		CobblerHost: cobblerHost,
	}

	scriptName := "postUpgrade.sh"
	scriptPath := filepath.Join(scriptDir, scriptName)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return "", fmt.Errorf("failed to generate %s", scriptName)
	}
	return scriptName, nil
}

// ReadContainerData returns values used to perform migration.
func ReadContainerData(scriptDir string) (string, string, string, error) {
	data, err := os.ReadFile(filepath.Join(scriptDir, "data"))
	if err != nil {
		return "", "", "", errors.New("failed to read data extracted from source host")
	}
	viper.SetConfigType("env")
	if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return "", "", "", fmt.Errorf("cannot read config: %s", err)
	}
	return viper.GetString("Timezone"), viper.GetString("old_pg_version"), viper.GetString("new_pg_version"), nil
}

// RunMigration execute the migration script.
func RunMigration(cnx *shared.Connection, tmpPath string, scriptName string) error {
	log.Info().Msg("Migrating server")
	err := ExecCommand(zerolog.InfoLevel, cnx, "/var/lib/uyuni-tools/"+scriptName)
	if err != nil {
		return fmt.Errorf("error running the migration script: %s", err)
	}
	return nil
}

// GenerateMigrationScript generates the script that perform migration.
func GenerateMigrationScript(sourceFqdn string, kubernetes bool) (string, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %s", err)
	}

	data := templates.MigrateScriptTemplateData{
		Volumes:    utils.ServerVolumeMounts,
		SourceFqdn: sourceFqdn,
		Kubernetes: kubernetes,
	}

	scriptPath := filepath.Join(scriptDir, "migrate.sh")
	if err = utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return "", fmt.Errorf("failed to generate migration script: %s", err)
	}

	return scriptDir, nil
}

// RunningImage returns the image running in the current system.
func RunningImage(cnx *shared.Connection, containerName string) (string, error) {
	command, err := cnx.GetCommand()

	switch command {
	case "podman":
		args := []string{"ps", "--format", "{{.Image}}", "--noheading"}
		image, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", args...)
		if err != nil {
			return "", err
		}
		return strings.Trim(string(image), "\n"), nil

	case "kubectl":

		//FIXME this will work until containers 0 is uyuni. Then jsonpath should be something like
		// {.items[0].spec.containers[?(@.name=="` + containerName + `")].image but there are problems
		// using RunCmdOutput with an arguments with round brackets
		args := []string{"get", "pods", kubernetes.ServerFilter, "-o", "jsonpath={.items[0].spec.containers[0].image}"}
		image, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", args...)

		log.Info().Msgf("image is: %s", image)
		if err != nil {
			return "", err
		}
		return strings.Trim(string(image), "\n"), nil
	}

	return command, err
}
