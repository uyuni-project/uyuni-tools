// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ExecCommand execute commands passed as argument in the current system.
func ExecCommand(logLevel zerolog.Level, cnx *shared.Connection, args ...string) error {
	podName, err := cnx.GetPodName()
	if err != nil {
		return utils.Errorf(err, L("exec command failed"))
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
func GenerateFinalizePostgresScript(
	scriptDir string, runAutotune bool, runReindex bool, runSchemaUpdate bool, migration bool, kubernetes bool,
) (string, error) {
	data := templates.FinalizePostgresTemplateData{
		RunAutotune:     runAutotune,
		RunReindex:      runReindex,
		RunSchemaUpdate: runSchemaUpdate,
		Migration:       migration,
		Kubernetes:      kubernetes,
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

// RunMigration execute the migration script.
func RunMigration(cnx *shared.Connection, tmpPath string, scriptName string) error {
	log.Info().Msg(L("Migrating server"))
	err := ExecCommand(zerolog.InfoLevel, cnx, "/var/lib/uyuni-tools/"+scriptName)
	if err != nil {
		return utils.Errorf(err, L("error running the migration script"))
	}
	return nil
}

// GenerateMigrationScript generates the script that perform migration.
func GenerateMigrationScript(sourceFqdn string, user string, kubernetes bool) (string, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	if err != nil {
		return "", utils.Errorf(err, L("failed to create temporary directory"))
	}

	data := templates.MigrateScriptTemplateData{
		Volumes:    utils.ServerVolumeMounts,
		SourceFqdn: sourceFqdn,
		User:       user,
		Kubernetes: kubernetes,
	}

	scriptPath := filepath.Join(scriptDir, "migrate.sh")
	if err = utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return "", utils.Errorf(err, L("failed to generate migration script"))
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

// SanityCheck verifies if an upgrade can be run.
func SanityCheck(cnx *shared.Connection, inspectedValues *utils.ServerInspectData, serverImage string) error {
	isUyuni, err := isUyuni(cnx)
	if err != nil {
		return utils.Errorf(err, L("cannot check server release"))
	}
	isCurrentUyuni := inspectedValues.UyuniRelease != ""
	isCurrentSuma := inspectedValues.SuseManagerRelease != ""

	if isUyuni && isCurrentSuma {
		return fmt.Errorf(
			L("currently SUSE Manager %s is installed, instead the image is Uyuni. Upgrade is not supported"),
			inspectedValues.SuseManagerRelease,
		)
	}

	if !isUyuni && isCurrentUyuni {
		return fmt.Errorf(
			L("currently Uyuni %s is installed, instead the image is SUSE Manager. Upgrade is not supported"),
			inspectedValues.UyuniRelease,
		)
	}

	if isUyuni {
		cnx_args := []string{"s/Uyuni release //g", "/etc/uyuni-release"}
		current_uyuni_release, err := cnx.Exec("sed", cnx_args...)
		if err != nil {
			return utils.Errorf(err, L("failed to read current uyuni release"))
		}
		log.Debug().Msgf("Current release is %s", string(current_uyuni_release))
		if inspectedValues.UyuniRelease != "" {
			return fmt.Errorf(L("cannot fetch release from image %s"), serverImage)
		}
		log.Debug().Msgf("Image %s is %s", serverImage, inspectedValues.UyuniRelease)
		if utils.CompareVersion(inspectedValues.UyuniRelease, string(current_uyuni_release)) < 0 {
			return fmt.Errorf(
				L("cannot downgrade from version %[1]s to %[2]s"),
				string(current_uyuni_release), inspectedValues.UyuniRelease,
			)
		}
	} else {
		cnx_args := []string{"s/SUSE Manager release //g", "/etc/susemanager-release"}
		current_suse_manager_release, err := cnx.Exec("sed", cnx_args...)
		if err != nil {
			return utils.Errorf(err, L("failed to read current susemanager release"))
		}
		log.Debug().Msgf("Current release is %s", string(current_suse_manager_release))
		if inspectedValues.SuseManagerRelease != "" {
			return fmt.Errorf(L("cannot fetch release from image %s"), serverImage)
		}
		log.Debug().Msgf("Image %s is %s", serverImage, inspectedValues.SuseManagerRelease)
		if utils.CompareVersion(inspectedValues.SuseManagerRelease, string(current_suse_manager_release)) < 0 {
			return fmt.Errorf(
				L("cannot downgrade from version %[1]s to %[2]s"),
				string(current_suse_manager_release), inspectedValues.SuseManagerRelease,
			)
		}
	}

	if inspectedValues.ImagePgVersion != "" {
		return fmt.Errorf(L("cannot fetch postgresql version from %s"), serverImage)
	}
	log.Debug().Msgf("Image %s has PostgreSQL %s", serverImage, inspectedValues.ImagePgVersion)
	if inspectedValues.CurrentPgVersion != "" {
		return fmt.Errorf(L("posgresql is not installed in the current deployment"))
	}
	log.Debug().Msgf("Current deployment has PostgreSQL %s", inspectedValues.CurrentPgVersion)

	return nil
}

func isUyuni(cnx *shared.Connection) (bool, error) {
	cnx_args := []string{"/etc/uyuni-release"}
	_, err := cnx.Exec("cat", cnx_args...)
	if err != nil {
		cnx_args := []string{"/etc/susemanager-release"}
		_, err := cnx.Exec("cat", cnx_args...)
		if err != nil {
			return false, errors.New(L("cannot find neither /etc/uyuni-release nor /etc/susemanagere-release"))
		}
		return false, nil
	}
	return true, nil
}
