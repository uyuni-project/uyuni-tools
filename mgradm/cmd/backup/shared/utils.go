// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	//lint:ignore ST1001 Ignore warning on lang tool import
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"golang.org/x/sys/unix"
)

const PodmanConfBackupFile = "podmanBackup.tar"
const SystemdConfBackupFile = "systemdBackup.tar"
const NetworkOutputFile = "uyuniNetwork.json"
const SecretBackupFile = "secrets.json"

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

func SanityChecks(outputDirectory string, dryRun bool) error {
	if utils.FileExists(outputDirectory) {
		if !utils.IsEmptyDirectory(outputDirectory) {
			return fmt.Errorf(L("output directory %s already exists and is not empty"), outputDirectory)
		}
	} else {
		if dryRun {
			log.Info().Msgf(L("Would create '%s' output directory"), outputDirectory)
		} else {
			if err := os.Mkdir(outputDirectory, 0622); err != nil {
				return fmt.Errorf(L("unable to create target output directory: %w"), err)
			}
		}
	}

	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}

	hostData, err := podman.InspectHost()
	if err != nil {
		return err
	}

	if !hostData.HasUyuniServer {
		return errors.New(L("server is not initialized."))
	}

	return nil
}
