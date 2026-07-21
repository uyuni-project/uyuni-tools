// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"os"
	"strings"

	"github.com/rs/zerolog"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// DefaultStorageRoot is the fallback path for podman storage when detection fails.
const DefaultStorageRoot = "/var/lib/containers/storage"

// IsCheckSkipped returns true if the user has set the UYUNI_SKIP_PREREQUISITES environment variable.
func IsCheckSkipped() bool {
	return os.Getenv("UYUNI_SKIP_PREREQUISITES") != ""
}

// CheckPrerequisites runs all pre-installation sanity checks.
func CheckPrerequisites(minMemoryGB, minStorageGB uint64, ports []types.PortMap) error {
	if IsCheckSkipped() {
		return nil
	}

	if err := utils.CheckMemory(minMemoryGB); err != nil {
		return err
	}

	storageRoot, err := GetPodmanVolumeBasePath()
	if err != nil || storageRoot == "" {
		storageRoot = DefaultStorageRoot
	}
	if err := utils.CheckStorage(storageRoot, minStorageGB); err != nil {
		return err
	}

	for _, portMap := range ports {
		if err := utils.CheckPort(portMap.Exposed); err != nil {
			return err
		}
	}

	return CheckPodmanRunningContainers()
}

// CheckPodmanRunningContainers checks if there are running containers on the uyuni network.
func CheckPodmanRunningContainers() error {
	out, err := runner("podman", "ps", "-q", "--filter", "network="+UyuniNetwork).
		Log(zerolog.DebugLevel).
		Exec()
	if err != nil {
		return utils.Errorf(err, L("failed to check running podman containers"))
	}

	if len(strings.TrimSpace(string(out))) > 0 {
		return errors.New(
			L("there are running containers on the uyuni network. Please stop them before installing a new one."),
		)
	}

	return nil
}
