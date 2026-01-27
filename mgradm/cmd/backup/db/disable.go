// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/rs/zerolog/log"

	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func Disable(flags *Flagpole) error {
	log.Info().Msg(L("Disable DB backup"))

	// Modify postgresql.conf and set archive_mode to off
	updates := map[string]string{
		"archive_mode": "off",
	}
	if err := UpdatePostgresConfig(updates); err != nil {
		return err
	}

	// Reload postgres config if postgres container is running
	cnx := shared.NewConnection("podman", podman.DBContainerName, "")
	if _, err := cnx.GetPodName(); err == nil {
		log.Debug().Msg("Reloading postgresql config")
		if _, err := cnx.Exec("/usr/bin/psql", "-U", "postgres", "-tAc", "SELECT pg_reload_conf();"); err != nil {
			return err
		}
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

	return nil
}
