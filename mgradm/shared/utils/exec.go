// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"os/exec"
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
		return utils.Error(err, L("exec command failed"))
	}

	commandArgs := []string{"exec", podName}

	command, err := cnx.GetCommand()
	if err != nil {
		log.Fatal().Err(err)
	}

	if command == "kubectl" {
		namespace, err := cnx.GetNamespace("")
		if namespace == "" {
			return utils.Error(err, L("failed retrieving namespace"))
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
	oldPgVersion string,
	newPgVersion string,
	backupDir string,
) (string, error) {
	data := templates.PostgreSQLVersionUpgradeTemplateData{
		OldVersion: oldPgVersion,
		NewVersion: newPgVersion,
		BackupDir:  backupDir,
	}

	scriptBuilder := new(strings.Builder)
	if err := data.Render(scriptBuilder); err != nil {
		return "", utils.Error(err, L("failed to render database upgrade script"))
	}
	return scriptBuilder.String(), nil
}

// GenerateFinalizePostgresScript generates the script to finalize PostgreSQL upgrade.
func GenerateFinalizePostgresScript(
	runReindex bool, runSchemaUpdate bool, migration bool, kubernetes bool,
) (string, error) {
	data := templates.FinalizePostgresTemplateData{
		RunReindex:      runReindex,
		RunSchemaUpdate: runSchemaUpdate,
		Migration:       migration,
		Kubernetes:      kubernetes,
	}

	scriptBuilder := new(strings.Builder)
	if err := data.Render(scriptBuilder); err != nil {
		return "", utils.Error(err, L("failed to render database finalization script"))
	}
	return scriptBuilder.String(), nil
}

// GeneratePostUpgradeScript generates the script to be run after upgrade.
func GeneratePostUpgradeScript() (string, error) {
	data := templates.PostUpgradeTemplateData{}

	scriptBuilder := new(strings.Builder)
	if err := data.Render(scriptBuilder); err != nil {
		return "", utils.Error(err, L("failed to render database post upgrade script"))
	}
	return scriptBuilder.String(), nil
}

// RunMigration execute the migration script.
func RunMigration(cnx *shared.Connection, scriptName string) error {
	log.Info().Msg(L("Migrating server"))
	err := ExecCommand(zerolog.InfoLevel, cnx, "/var/lib/uyuni-tools/"+scriptName)
	if err != nil {
		return utils.Error(err, L("error running the migration script"))
	}
	return nil
}

// GenerateMigrationScript generates the script that perform migration.
func GenerateMigrationScript(
	sourceFqdn string,
	user string,
	kubernetes bool,
	prepare bool,
	dbHost string,
	reportDBHost string,
) (string, error) {
	// For podman we want to backup tls certificates to the temporary volume we
	// later use when creating secrets.
	volumes := append(utils.ServerVolumeMounts, utils.VarPgsqlDataVolumeMount)
	if !kubernetes {
		volumes = append(volumes, utils.EtcTLSTmpVolumeMount)
	}

	data := templates.MigrateScriptTemplateData{
		Volumes:      volumes,
		SourceFqdn:   sourceFqdn,
		User:         user,
		Kubernetes:   kubernetes,
		Prepare:      prepare,
		DBHost:       dbHost,
		ReportDBHost: reportDBHost,
	}

	scriptBuilder := new(strings.Builder)
	if err := data.Render(scriptBuilder); err != nil {
		return "", utils.Error(err, L("failed to generate migration script"))
	}

	return scriptBuilder.String(), nil
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
	return nil
}
