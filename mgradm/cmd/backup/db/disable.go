// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func Disable(flags *Flagpole) error {
	log.Info().Msg(L("Disabling DB backup"))

	wasRunning := systemd.IsServiceRunning(podman.DBService)
	cnx := shared.NewUserConnection("podman", podman.DBContainerName, "", podman.DBRuntimeUser)
	if wasRunning {
		if !flags.Force {
			res, err := utils.YesNo(L("Database service needs to be restarted to reload backup configurations, continue?"))
			if err != nil || !res {
				log.Info().Msg(L("Backup reconfiguration aborted"))
				return nil
			}
		}
		if err := systemd.StopService(podman.DBService); err != nil {
			return err
		}
	}

	// Modify postgresql.conf and set archive_mode to off
	updates := map[string]string{
		"archive_mode": "off",
	}
	if err := UpdatePostgresConfig(updates); err != nil {
		return err
	}

	if flags.Purge.Volume {
		if !flags.Force {
			res, err := utils.YesNo(L("Backup volume will be removed, continue?"))
			if err != nil || !res {
				log.Info().Msg(L("Aborting volume purge"))
				return nil
			}
		}
		if err := podman.DeleteVolume(utils.VarPgsqlBackupVolumeMount.Name, false); err != nil {
			return err
		}
		log.Info().Msgf(L("Backup volume '%s' removed"), utils.VarPgsqlBackupVolumeMount.Name)
	}

	if wasRunning {
		log.Info().Msg(L("Starting database…"))
		startAttempt := time.Now()
		if err := systemd.StartService(podman.DBService); err != nil {
			return err
		}
		if err := cnx.WaitForHealthcheck(startAttempt); err != nil {
			podman.ShowServiceLogs(podman.DBService, startAttempt)
			return err
		}
	}
	return Status()
}
