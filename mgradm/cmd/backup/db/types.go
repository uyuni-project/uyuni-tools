// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"errors"
	"fmt"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type Flagpole struct {
	Force bool `mapstructure:"force"`
	Purge struct {
		Volume bool `mapstructure:"volume"`
	} `mapstructure:"purge"`
}

var (
	ErrArchiveModeOff              = errors.New(L("archive_mode is off"))
	ErrArchiveCommandMisconfigured = errors.New(L("archive_command is misconfigured"))
	ErrArchiveMountMisconfigured   = errors.New(L("archive volume is not configured correctly"))
)

// For unit testing, enable overridable systemd.
var systemd podman.Systemd = podman.NewSystemd()

func ArchiveCommand() string {
	return fmt.Sprintf("'/usr/bin/smdba-pgarchive --source \"%%p\" --destination \"%s/%%f\"'",
		utils.VarPgsqlBackupVolumeMount.MountPath)
}

func RestoreCommand() string {
	return fmt.Sprintf("/usr/bin/cp %s/%%f %%p", utils.VarPgsqlBackupVolumeMount.MountPath)
}
