// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/exec"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var newExecCmd = exec.NewCommand

// NewCommand returns a new cobra.Command for term.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "term",
		Short: "Run a terminal inside the server container",
		RunE: func(cmd *cobra.Command, args []string) error {
			execCmd := newExecCmd(globalFlags)
			execArgs := []string{"-i", "-t"}
			backend, err := cmd.Flags().GetString("backend")
			if err == nil {
				execArgs = append(execArgs, "--backend", backend)
			}
			if err := execCmd.Flags().Parse(execArgs); err != nil {
				return err
			}
			return execCmd.RunE(execCmd, []string{"bash"})
		},
	}

	utils.AddBackendFlag(cmd)
	return cmd
}
