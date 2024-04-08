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
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// InspectScriptFilename is the inspect script basename.
var InspectScriptFilename = "inspect.sh"

var inspectValues = []types.InspectData{
	types.NewInspectData("uyuni_release", "cat /etc/*release | grep 'Uyuni release' | cut -d ' ' -f3 || true"),
	types.NewInspectData("suse_manager_release", "cat /etc/*release | grep 'SUSE Manager release' | cut -d ' ' -f4 || true"),
	types.NewInspectData("fqdn", "cat /etc/rhn/rhn.conf 2>/dev/null | grep 'java.hostname' | cut -d' ' -f3 || true"),
	types.NewInspectData("image_pg_version", "rpm -qa --qf '%{VERSION}\\n' 'name=postgresql[0-8][0-9]-server'  | cut -d. -f1 | sort -n | tail -1 || true"),
	types.NewInspectData("current_pg_version", "(test -e /var/lib/pgsql/data/PG_VERSION && cat /var/lib/pgsql/data/PG_VERSION) || true"),
	types.NewInspectData("registration_info", "transactional-update --quiet register --status 2>/dev/null || true"),
	types.NewInspectData("scc_username", "(test -e /etc/zypp/credentials.d/SCCcredentials && cat /etc/zypp/credentials.d/SCCcredentials | grep username | cut -d= -f2) || true"),
	types.NewInspectData("scc_password", "(test -e /etc/zypp/credentials.d/SCCcredentials && cat /etc/zypp/credentials.d/SCCcredentials | grep password | cut -d= -f2) || true"),
}

// InspectOutputFile represents the directory and the basename where the inspect values are stored.
var InspectOutputFile = types.InspectFile{
	Directory: "/var/lib/uyuni-tools",
	Basename:  "data",
}

// ExecCommand execute commands passed as argument in the current system.
func ExecCommand(logLevel zerolog.Level, cnx *shared.Connection, args ...string) error {
	podName, err := cnx.GetPodName()
	if err != nil {
		return fmt.Errorf(L("exec command failed: %s"), err)
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
		return "", fmt.Errorf(L("failed to generate %s"), scriptName)
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
		return "", fmt.Errorf(L("failed to generate %s"), scriptName)
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
		return "", fmt.Errorf(L("failed to generate %s"), scriptName)
	}
	return scriptName, nil
}

// ReadContainerData returns values used to perform migration.
func ReadContainerData(scriptDir string) (string, string, string, error) {
	data, err := os.ReadFile(filepath.Join(scriptDir, "data"))
	if err != nil {
		return "", "", "", errors.New(L("failed to read data extracted from source host"))
	}
	viper.SetConfigType("env")
	if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return "", "", "", fmt.Errorf(L("cannot read config: %s"), err)
	}
	if len(viper.GetString("Timezone")) <= 0 {
		return "", "", "", errors.New(L("cannot retrieve timezone"))
	}
	if len(viper.GetString("old_pg_version")) <= 0 {
		return "", "", "", errors.New(L("cannot retrieve source PostgreSQL version"))
	}
	if len(viper.GetString("new_pg_version")) <= 0 {
		return "", "", "", errors.New(L("cannot retrieve image PostgreSQL version"))
	}
	return viper.GetString("Timezone"), viper.GetString("old_pg_version"), viper.GetString("new_pg_version"), nil
}

// RunMigration execute the migration script.
func RunMigration(cnx *shared.Connection, tmpPath string, scriptName string) error {
	log.Info().Msg(L("Migrating server"))
	err := ExecCommand(zerolog.InfoLevel, cnx, "/var/lib/uyuni-tools/"+scriptName)
	if err != nil {
		return fmt.Errorf(L("error running the migration script: %s"), err)
	}
	return nil
}

// GenerateMigrationScript generates the script that perform migration.
func GenerateMigrationScript(sourceFqdn string, user string, kubernetes bool) (string, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	if err != nil {
		return "", fmt.Errorf(L("failed to create temporary directory: %s"), err)
	}

	data := templates.MigrateScriptTemplateData{
		Volumes:    utils.ServerVolumeMounts,
		SourceFqdn: sourceFqdn,
		User:       user,
		Kubernetes: kubernetes,
	}

	scriptPath := filepath.Join(scriptDir, "migrate.sh")
	if err = utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return "", fmt.Errorf(L("failed to generate migration script: %s"), err)
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

		log.Info().Msgf(L("Image is: %s"), image)
		if err != nil {
			return "", err
		}
		return strings.Trim(string(image), "\n"), nil
	}

	return command, err
}

// ReadInspectData returns a map with the values inspected by an image and deploy.
func ReadInspectData(scriptDir string, prefix ...string) (map[string]string, error) {
	path := filepath.Join(scriptDir, "data")
	log.Debug().Msgf("Trying to read %s", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]string{}, fmt.Errorf(L("cannot parse file %s: %s"), path, err)
	}

	inspectResult := make(map[string]string)

	viper.SetConfigType("env")
	if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return map[string]string{}, fmt.Errorf(L("cannot read config: %s"), err)
	}

	for _, v := range inspectValues {
		if len(viper.GetString(v.Variable)) > 0 {
			index := v.Variable
			/* Just the first value of prefix is used.
			 * This slice is just to allow an empty argument
			 */
			if len(prefix) >= 1 {
				index = prefix[0] + v.Variable
			}
			inspectResult[index] = viper.GetString(v.Variable)
		}
	}
	return inspectResult, nil
}

// InspectHost check values on a host machine.
func InspectHost() (map[string]string, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return map[string]string{}, fmt.Errorf(L("failed to create temporary directory: %s"), err)
	}

	if err := GenerateInspectHostScript(scriptDir); err != nil {
		return map[string]string{}, err
	}

	if err := utils.RunCmdStdMapping(zerolog.DebugLevel, scriptDir+"/inspect.sh"); err != nil {
		return map[string]string{}, fmt.Errorf(L("failed to run inspect script in host system: %s"), err)
	}

	inspectResult, err := ReadInspectData(scriptDir, "host_")
	if err != nil {
		return map[string]string{}, fmt.Errorf(L("cannot inspect host data: %s"), err)
	}

	return inspectResult, err
}

// GenerateInspectContainerScript create the host inspect script.
func GenerateInspectHostScript(scriptDir string) error {
	data := templates.InspectTemplateData{
		Param:      inspectValues,
		OutputFile: scriptDir + "/" + InspectOutputFile.Basename,
	}

	scriptPath := filepath.Join(scriptDir, InspectScriptFilename)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return fmt.Errorf(L("failed to generate inspect script: %s"), err)
	}
	return nil
}

// GenerateInspectContainerScript create the container inspect script.
func GenerateInspectContainerScript(scriptDir string) error {
	data := templates.InspectTemplateData{
		Param:      inspectValues,
		OutputFile: InspectOutputFile.Directory + "/" + InspectOutputFile.Basename,
	}

	scriptPath := filepath.Join(scriptDir, InspectScriptFilename)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return fmt.Errorf(L("failed to generate inspect script: %s"), err)
	}
	return nil
}
