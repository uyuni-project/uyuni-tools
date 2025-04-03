// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"golang.org/x/sys/unix"
)

const PodmanConfBackupFile = "podmanBackup.tar"
const SystemdConfBackupFile = "systemdBackup.tar"
const NetworkOutputFile = "uyuniNetwork.json"
const SecretBackupFile = "secrets.json"

const VolumesSubdir = "volumes"
const ImagesSubdir = "images"

func StorageCheck(volumes []string, images []string, outputDirectory string) error {
	// check disk space availability based on volume work list and container image list
	var outStat unix.Statfs_t
	if err := unix.Statfs(outputDirectory, &outStat); err != nil {
		log.Warn().Err(err).Msgf(L("unable to determine target %s storage size"), outputDirectory)
	}
	freeSpace := outStat.Bavail * uint64(outStat.Bsize)
	var spaceRequired int64

	// calculate required space
	for _, volume := range volumes {
		mountPoint, err := podman.GetVolumeMountPoint(volume)
		if err != nil {
			return err
		}
		volumeSize, err := dirSize(mountPoint)
		if err != nil {
			return err
		}
		spaceRequired += volumeSize
	}

	// Calculate the size of the images
	for _, image := range images {
		// This is over estimating the actual size on disk since the layers can be shared,
		// but that can't be bad to have more disk than actually needed.
		size, err := podman.GetImageVirtualSize(image)
		if err != nil {
			return err
		}
		spaceRequired += size
	}

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

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		size += info.Size()
		return nil
	})
	return size, err
}
