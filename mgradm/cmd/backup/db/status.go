// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/pgsql"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func Status() error {
	log.Info().Msg(L("Checking DB backup status"))

	status := "enabled"
	var result error
	err := CheckStatus()
	if err != nil {
		if errors.Is(err, ErrArchiveModeOff) {
			status = "disabled"
		} else if errors.Is(err, ErrArchiveCommandMisconfigured) || errors.Is(err, ErrArchiveMountMisconfigured) {
			status = "misconfigured"
		} else {
			status = "unknown"
			result = err
		}
	}
	if zerolog.GlobalLevel() < zerolog.WarnLevel {
		log.Info().Err(err).Msgf(L("Database continuous backup is %[1]s. Backup volume is %[2]s"), status,
			utils.VarPgsqlBackupVolumeMount.Name)
	} else {
		// Assuming this is called by a script when logLevel is higher then info
		fmt.Println(status)
	}
	return result
}

// Check status of the database wal backup.
// Reports nil if backup is enabled and correctly configured. Otherwise reports an error.
func CheckStatus() error {
	log.Debug().Msg("Checking for the database configuration")
	cnx := shared.NewConnection("podman", podman.DBContainerName, "")
	if _, err := cnx.GetPodName(); err == nil {
		log.Debug().Msg("Checking for runtime config")
		if err := checkStatusSQL(cnx); err != nil {
			return err
		}
	}
	log.Debug().Msg("Checking for file based config")
	if err := checkStatusFile(); err != nil {
		return err
	}

	// If we are here, configuration is set for archiving
	log.Debug().Msg("Checking for backup mount presence")
	if err := CheckBackupMountPresence(); err != nil {
		return err
	}
	return nil
}

func checkStatusSQL(cnx *shared.Connection) error {
	out, err := cnx.Exec("/usr/bin/psql", "-U", "postgres", "-tAc", "SHOW archive_mode;")
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(out)) != "on" {
		return ErrArchiveModeOff
	}

	out, err = cnx.Exec("/usr/bin/psql", "-U", "postgres", "-tAc", "SHOW archive_command;")
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(out)) != strings.Trim(ArchiveCommand(), "'") {
		return ErrArchiveCommandMisconfigured
	}
	return nil
}

func checkStatusFile() error {
	config, err := ParsePostgresConfig()
	if err != nil {
		return err
	}

	if val, ok := config["archive_mode"]; !ok || val != "on" {
		return ErrArchiveModeOff
	}
	if val, ok := config["archive_command"]; !ok || val != ArchiveCommand() {
		log.Trace().Msgf("archive_command: \"%s\"", config["archive_command"])
		log.Trace().Msgf("expectd_command: \"%s\"", ArchiveCommand())
		return ErrArchiveCommandMisconfigured
	}
	return nil
}

func CheckBackupMountPresence() error {
	dropinPaths, err := systemd.GetServiceProperty(podman.DBService, podman.DropInPaths)
	if err != nil {
		return err
	}

	if !podman.IsVolumePresent(utils.VarPgsqlBackupVolumeMount.Name) {
		log.Debug().Msg("backup volume not mounted")
		return ErrArchiveMountMisconfigured
	}

	confFileFound := false
	for _, confPath := range strings.Split(strings.TrimPrefix(dropinPaths, "DropInPaths="), " ") {
		if strings.HasSuffix(confPath, pgsql.BackupVolumeConfigName) {
			content, err := os.ReadFile(confPath)
			if err != nil {
				return err
			}
			if strings.Contains(string(content), pgsql.BackupVolumeConfig()) {
				confFileFound = true
			}
			break
		}
	}

	if !confFileFound {
		log.Debug().Msg("valid backup config not found")
		return ErrArchiveMountMisconfigured
	}
	return nil
}
