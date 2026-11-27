// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgradm/cmd"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestSubCommandsHaveGroup(t *testing.T) {
	mgradmCmd, _ := cmd.NewUyuniadmCommand()
	testutils.AssertTrue(t, "There's at least one mgradm subcommand without group", mgradmCmd.AllChildCommandsHaveGroup())
}
