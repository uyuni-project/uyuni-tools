// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgradm/cmd"
)

func TestSubCommandsHaveGroup(t *testing.T) {
	mgradmCmd, _ := cmd.NewUyuniadmCommand()
	if !mgradmCmd.AllChildCommandsHaveGroup() {
		t.Errorf("There's at least one mgradm subcommand without group")
	}
}
