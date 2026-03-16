// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
)

func Rebase() error {
	log.Info().Msg(L("Rebasing postgres WAL backup to new basebackup"))

	if err := CheckStatus(); err != nil {
		return fmt.Errorf(L("database backup is not correctly configured: %w"), err)
	}

	cnx := shared.NewConnection("podman", podman.DBContainerName, "")
	isRunning, _ := cnx.GetPodName()
	wasRunning := isRunning != ""

	if !wasRunning {
		log.Info().Msg(L("Starting database..."))
		if err := systemd.StartService(podman.DBService); err != nil {
			return err
		}
	}
	if err := cnx.WaitForHealthcheck(); err != nil {
		return err
	}

	if err := RunBaseBackup(cnx); err != nil {
		return err
	}

	if !wasRunning {
		log.Info().Msg(L("Stopping database..."))
		if err := systemd.StopService(podman.DBService); err != nil {
			log.Warn().Err(err).Send()
		}
	}

	return nil
}
