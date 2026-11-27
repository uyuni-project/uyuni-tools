// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/exec"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// Ensure the term command properly delegates to the exec one.
func TestExecute(t *testing.T) {
	var globalFlags types.GlobalFlags

	newExecCmd = func(globalFlags *types.GlobalFlags) *cobra.Command {
		execCmd := exec.NewCommand(globalFlags)
		execCmd.RunE = func(cmd *cobra.Command, _ []string) error {
			interactive, err := cmd.Flags().GetBool("interactive")
			testutils.AssertNoError(t, "failed to get interactive flag", err)
			testutils.AssertTrue(t, "interactive flag not passed", interactive)

			tty, err := cmd.Flags().GetBool("tty")
			testutils.AssertNoError(t, "failed to get tty flag", err)
			testutils.AssertTrue(t, "tty flag not passed", tty)

			backend, err := cmd.Flags().GetString("backend")
			testutils.AssertNoError(t, "failed to get backend flag", err)
			testutils.AssertEquals(t, "backend flag not passed", "mybackend", backend)

			return errors.New("some error")
		}
		return execCmd
	}

	cmd := NewCommand(&globalFlags)
	testutils.AssertNoError(t, "failed to parse flags", cmd.Flags().Parse([]string{"--backend", "mybackend"}))
	testutils.AssertError(t, "some error", cmd.RunE(cmd, []string{}))
}
