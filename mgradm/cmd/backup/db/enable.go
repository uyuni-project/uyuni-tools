// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/pgsql"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func Enable(force bool) error {
	log.Info().Msg(L("Enable DB backup"))

	if err := CheckStatus(); err == nil && !force {
		return errors.New(L("backup is already configured. Use --force to reconfigure"))
	}

	cnx := shared.NewConnection("podman", podman.DBContainerName, "")
	isRunning, _ := cnx.GetPodName()
	wasRunning := isRunning != ""

	if wasRunning {
		if !force {
			res, err := utils.YesNo(L("Database service needs to be restarted to reload backup configurations, continue"))
			if err != nil || !res {
				log.Info().Msg(L("Backup configuration aborted"))
				return nil
			}
		}
		if err := systemd.StopService(podman.DBService); err != nil {
			return err
		}
	}

	// We need to also regenerate service
	if err := pgsql.GenerateBackupVolumeConfig(systemd); err != nil {
		return err
	}

	// Modify postgresql.conf and set archive_command to ArchiveCommand and archive_mode to yes,
	updates := map[string]string{
		"archive_mode":    "on",
		"archive_command": ArchiveCommand(),
		"wal_level":       "replica",
	}
	if err := UpdatePostgresConfig(updates); err != nil {
		return err
	}

	// Check if all is well configured
	if err := CheckStatus(); err != nil {
		return err
	}

	// Prepare initial snapshot. For this container either must be running or we need to start container database
	log.Info().Msg(L("Starting database..."))
	if err := systemd.StartService(podman.DBService); err != nil {
		return err
	}
	if err := cnx.WaitForHealthcheck(); err != nil {
		return err
	}

	data := templates.EnablePostgresTemplateData{
		BackupDir: utils.VarPgsqlBackupVolumeMount.MountPath,
	}

	scriptBuilder := new(strings.Builder)
	if err := data.Render(scriptBuilder); err != nil {
		return utils.Error(err, L("failed to generate postgresql backup script"))
	}

	// We need to exec enable script inside the database container
	if _, err := cnx.ExecScript(scriptBuilder.String()); err != nil {
		return err
	}

	// Only stop the database if it was not running at the beginning
	if !wasRunning {
		if err := systemd.StopService(podman.DBService); err != nil {
			log.Warn().Err(err).Send()
		}
	}

	log.Info().Msgf(L("Continuous Archiving backup configured. Backup target volume is '%s'"),
		utils.VarPgsqlBackupVolumeMount.Name)

	return nil
}
