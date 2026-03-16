// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rs/zerolog/log"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func Restore(force bool) error {
	log.Info().Msg(L("Restoring DB backup"))

	if !force {
		if res, err := utils.YesNo(L("Restoring from backup is a destructive operation. Proceed")); err != nil || !res {
			log.Info().Msg(L("Aborting"))
			return nil
		}
	}
	log.Info().Msg(L("Stopping database service..."))
	if err := systemd.StopService(podman.DBService); err != nil {
		return err
	}

	image := podman.GetServiceImage(podman.DBService)
	if image == "" {
		return errors.New(L("failed to determine database image"))
	}

	volumes := []types.VolumeMount{
		utils.VarPgsqlDataVolumeMount,
		utils.VarPgsqlBackupVolumeMount,
	}

	// Actual data moving is in the restore script rendered and executed below
	data := templates.RestorePostgresTemplateData{
		Datadir:    utils.VarPgsqlDataVolumeMount.MountPath,
		Basebackup: path.Join(utils.VarPgsqlBackupVolumeMount.MountPath, "base.tar.gz"),
	}

	scriptBuilder := new(strings.Builder)
	if err := data.Render(scriptBuilder); err != nil {
		return utils.Error(err, L("failed to generate postgresql restore script"))
	}

	log.Info().Msg(L("Restoring base backup..."))
	if err := podman.RunContainer("uyuni-restore", image, volumes, []string{},
		[]string{"bash", "-e", "-c", scriptBuilder.String()}); err != nil {
		return err
	}

	// Modify postgresql.conf and set restore_command to RestoreCommand
	updates := map[string]string{
		"restore_command": RestoreCommand(),
	}
	if err := UpdatePostgresConfig(updates); err != nil {
		return err
	}

	log.Info().Msg(L("Starting database service..."))
	if err := systemd.StartService(podman.DBService); err != nil {
		return err
	}

	log.Info().Msg(L("Base backup restore complete. Database is recovering."))

	mountPoint, err := podman.GetVolumeMountPoint(utils.VarPgsqlDataVolumeMount.Name)
	if err != nil {
		return err
	}
	recoverySignalPath := path.Join(mountPoint, "recovery.signal")

	if _, err := os.Stat(recoverySignalPath); err != nil {
		if os.IsNotExist(err) {
			log.Info().Msg(L("Database is restored."))
			return nil
		}
		return utils.Error(err, L("failed to check database recovery status"))
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = L(" Database is recovering... (use Ctrl-C to interrupt waiting)")
	s.Start()
	for {
		if _, err := os.Stat(recoverySignalPath); err != nil {
			s.Stop()
			if os.IsNotExist(err) {
				log.Info().Msg(L("Database is restored."))
				return nil
			}
			return utils.Error(err, L("error while waiting for database recovery to complete, check database logs"))
		}
		time.Sleep(2 * time.Second)
	}
}
