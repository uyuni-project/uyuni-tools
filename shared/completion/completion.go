// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package completion

import (
	"os"

	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// NewCommand  command for generates completion script.
func NewCommand(_ *types.GlobalFlags) *cobra.Command {
	shellCompletionCmd := &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 L("Generate shell completion script"),
		Long:                  L("Generate shell completion script"),
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Hidden:                true,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				if err := cmd.Root().GenBashCompletion(os.Stdout); err != nil {
					return utils.Errorf(err, L("cannot generate %s completion"), args[0])
				}
			case "zsh":
				if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
					return utils.Errorf(err, L("cannot generate %s completion"), args[0])
				}
			case "fish":
				if err := cmd.Root().GenFishCompletion(os.Stdout, true); err != nil {
					return utils.Errorf(err, L("cannot generate %s completion"), args[0])
				}
			}
			return nil
		},
	}
	return shellCompletionCmd
}
