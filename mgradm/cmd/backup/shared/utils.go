// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"errors"
	"os/exec"

	//lint:ignore ST1001 Ignore warning on lang tool import

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"golang.org/x/sys/unix"
)

const PodmanConfBackupFile = "podmanBackup.tar"
const SystemdConfBackupFile = "systemdBackup.tar"
const NetworkOutputFile = "uyuniNetwork.json"
const SecretBackupFile = "secrets.json"

const VolumesSubdir = "volumes"
const ImagesSubdir = "images"

// runCmd* are a function pointers to use for easies unit testing.
var runCmdOutput = utils.RunCmdOutput

func StorageCheck(volumes []string, images []string, outputDirectory string) error {
	// check disk space availability based on volume work list and container image list
	var outStat unix.Statfs_t
	if err := unix.Statfs(outputDirectory, &outStat); err != nil {
		log.Warn().Err(err).Msgf(L("unable to determine target %s storage size"), outputDirectory)
	}
	freeSpace := outStat.Bavail * uint64(outStat.Bsize)
	spaceRequired := 0
	// TODO calculate required space

	if freeSpace < uint64(spaceRequired) {
		return errors.New(L("insufficient space on target device"))
	}
	return nil
}

func SanityChecks() error {
	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}

	return nil
}
