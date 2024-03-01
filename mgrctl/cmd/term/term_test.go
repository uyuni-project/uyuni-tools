// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/exec"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// Ensure the term command properly delegates to the exec one.
func TestExecute(t *testing.T) {
	var globalFlags types.GlobalFlags

	newExecCmd = func(globalFlags *types.GlobalFlags) *cobra.Command {
		execCmd := exec.NewCommand(globalFlags)
		execCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if interactive, err := cmd.Flags().GetBool("interactive"); err != nil || !interactive {
				t.Error("interactive flag not passed")
			}
			if tty, err := cmd.Flags().GetBool("tty"); err != nil || !tty {
				t.Error("tty flag not passed")
			}
			if backend, err := cmd.Flags().GetString("backend"); err != nil || backend != "mybackend" {
				t.Error("backend flag not passed")
			}
			return errors.New("some error")
		}
		return execCmd
	}

	cmd := NewCommand(&globalFlags)
	if err := cmd.Flags().Parse([]string{"--backend", "mybackend"}); err != nil {
		t.Errorf("failed to parse flags: %s", err)
	}
	if err := cmd.RunE(cmd, []string{}); err.Error() != "some error" {
		t.Errorf("Unexpected error returned")
	}
}
