// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared"
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
		namespace, err := cnx.GetNamespace("")
		if namespace == "" {
			return utils.Errorf(err, L("failed retrieving namespace"))
		}
		commandArgs = append(commandArgs, "-n", namespace, "-c", "uyuni", "--")
	}

	commandArgs = append(commandArgs, "sh", "-c", strings.Join(args, " "))

	runCmd := exec.Command(command, commandArgs...)
	logger := log.Logger.Level(logLevel)
	runCmd.Stdout = logger
	runCmd.Stderr = logger
	return runCmd.Run()
}

// GeneratePgsqlVersionUpgradeScript generates the PostgreSQL version upgrade script.
func GeneratePgsqlVersionUpgradeScript(
	scriptDir string,
	oldPgVersion string,
	newPgVersion string,
	kubernetes bool,
) (string, error) {
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
func GeneratePostUpgradeScript(scriptDir string) (string, error) {
	data := templates.PostUpgradeTemplateData{}

	scriptName := "postUpgrade.sh"
	scriptPath := filepath.Join(scriptDir, scriptName)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return "", fmt.Errorf(L("failed to generate %s"), scriptName)
	}
	return scriptName, nil
}

// RunMigration execute the migration script.
func RunMigration(cnx *shared.Connection, scriptName string) error {
	log.Info().Msg(L("Migrating server"))
	err := ExecCommand(zerolog.InfoLevel, cnx, "/var/lib/uyuni-tools/"+scriptName)
	if err != nil {
		return utils.Errorf(err, L("error running the migration script"))
	}
	return nil
}

// GenerateMigrationScript generates the script that perform migration.
func GenerateMigrationScript(sourceFqdn string, user string, kubernetes bool, prepare bool) (string, func(), error) {
	scriptDir, cleaner, err := utils.TempDir()
	if err != nil {
		return "", nil, err
	}

	data := templates.MigrateScriptTemplateData{
		Volumes:    utils.ServerVolumeMounts,
		SourceFqdn: sourceFqdn,
		User:       user,
		Kubernetes: kubernetes,
		Prepare:    prepare,
	}

	scriptPath := filepath.Join(scriptDir, "migrate.sh")
	if err = utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return "", cleaner, utils.Errorf(err, L("failed to generate migration script"))
	}

	return scriptDir, cleaner, nil
}

// SanityCheck verifies if an upgrade can be run.
func SanityCheck(
	runningValues *utils.ServerInspectData,
	inspectedValues *utils.ServerInspectData,
	serverImage string,
) error {
	// Skip the uyuni / SUSE Manager release checks if the runningValues is nil.
	if runningValues != nil {
		isUyuni := runningValues.UyuniRelease != ""
		isUyuniImage := inspectedValues.UyuniRelease != ""
		isSumaImage := inspectedValues.SuseManagerRelease != ""

		if isUyuni && isSumaImage {
			return fmt.Errorf(
				L("currently SUSE Manager %s is installed, instead the image is Uyuni. Upgrade is not supported"),
				inspectedValues.SuseManagerRelease,
			)
		}

		if !isUyuni && isUyuniImage {
			return fmt.Errorf(
				L("currently Uyuni %s is installed, instead the image is SUSE Manager. Upgrade is not supported"),
				inspectedValues.UyuniRelease,
			)
		}

		if isUyuni {
			currentUyuniRelease := runningValues.UyuniRelease
			log.Debug().Msgf("Current release is %s", string(currentUyuniRelease))
			if !isUyuniImage {
				return fmt.Errorf(L("cannot fetch release from image %s"), serverImage)
			}
			log.Debug().Msgf("Image %s is %s", serverImage, inspectedValues.UyuniRelease)
			if utils.CompareVersion(inspectedValues.UyuniRelease, string(currentUyuniRelease)) < 0 {
				return fmt.Errorf(
					L("cannot downgrade from version %[1]s to %[2]s"),
					string(currentUyuniRelease), inspectedValues.UyuniRelease,
				)
			}
		} else {
			currentSuseManagerRelease := runningValues.SuseManagerRelease
			log.Debug().Msgf("Current release is %s", currentSuseManagerRelease)
			if !isSumaImage {
				return fmt.Errorf(L("cannot fetch release from image %s"), serverImage)
			}
			log.Debug().Msgf("Image %s is %s", serverImage, inspectedValues.SuseManagerRelease)
			if utils.CompareVersion(inspectedValues.SuseManagerRelease, currentSuseManagerRelease) < 0 {
				return fmt.Errorf(
					L("cannot downgrade from version %[1]s to %[2]s"),
					currentSuseManagerRelease, inspectedValues.SuseManagerRelease,
				)
			}
		}
	}

	// Perform PostgreSQL version checks.
	if inspectedValues.ImagePgVersion == "" {
		return fmt.Errorf(L("cannot fetch PostgreSQL version from %s"), serverImage)
	}
	log.Debug().Msgf("Image %s has PostgreSQL %s", serverImage, inspectedValues.ImagePgVersion)
	if inspectedValues.CurrentPgVersion == "" {
		return errors.New(L("PostgreSQL is not installed in the current deployment"))
	}
	log.Debug().Msgf("Current deployment has PostgreSQL %s", inspectedValues.CurrentPgVersion)

	return nil
}
