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

const skipPrerequisitesEnv = "UYUNI_SKIP_PREREQUISITES"

// IsCheckSkipped returns whether prerequisite checks have been explicitly disabled.
func IsCheckSkipped() bool {
	value := os.Getenv(skipPrerequisitesEnv)
	return value == "1" || strings.EqualFold(value, "true")
}

// CheckPrerequisites runs all pre-installation sanity checks.
func CheckPrerequisites(minMemoryGB, minStorageGB uint64, ports []types.PortMap, network string) error {
	if IsCheckSkipped() {
		return nil
	}

	storageRoot, err := GetPodmanVolumeBasePath()
	if err != nil {
		return err
	}

	errs := utils.JoinErrors(
		utils.CheckMemory(minMemoryGB),
		utils.CheckStorage(storageRoot, minStorageGB),
		CheckPodmanRunningContainers(network),
	)
	for _, portMap := range ports {
		errs = utils.JoinErrors(errs, utils.CheckPort(portMap))
	}
	return errs
}

// CheckPodmanRunningContainers checks if there are running containers on the given network.
func CheckPodmanRunningContainers(network string) error {
	out, err := runner("podman", "ps", "-q", "--filter", "network="+network).
		Log(zerolog.DebugLevel).
		Exec()
	if err != nil {
		return utils.Errorf(err, L("failed to check running podman containers"))
	}

	if len(strings.TrimSpace(string(out))) > 0 {
		return errors.New(
			L("there are running containers on the Uyuni network. Please stop them before installing."),
		)
	}

	return nil
}
