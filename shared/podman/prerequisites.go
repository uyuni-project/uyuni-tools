// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"strings"

	"github.com/rs/zerolog"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// CheckPodmanRunningContainers checks if there are running containers on the uyuni network.
func CheckPodmanRunningContainers() error {
	out, err := runner("podman", "ps", "-q", "--filter", "network="+UyuniNetwork).
		Log(zerolog.DebugLevel).
		Exec()
	if err != nil {
		return utils.Errorf(err, L("failed to check running podman containers"))
	}

	if len(strings.TrimSpace(string(out))) > 0 {
		return errors.New(L("there are running containers on the uyuni network. Please stop them before installing or upgrading (see issue #323)."))
	}

	return nil
}
