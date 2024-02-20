// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package completion

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand  command for generates completion script.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	shellCompletionCmd := &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate shell completion script",
		Long:                  "Generate shell completion script",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish"},
		Args:                  cobra.ExactValidArgs(1),
		Hidden:                true,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				if err := cmd.Root().GenBashCompletion(os.Stdout); err != nil {
					return fmt.Errorf("cannot generate bash completion: %s", err)
				}
			case "zsh":
				if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
					return fmt.Errorf("cannot generate zsh completion: %s", err)
				}
			case "fish":
				if err := cmd.Root().GenFishCompletion(os.Stdout, true); err != nil {
					return fmt.Errorf("cannot generate fish completion: %s", err)
				}
			}
			return nil
		},
	}
	return shellCompletionCmd
}
