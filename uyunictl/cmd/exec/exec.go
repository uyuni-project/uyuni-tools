package exec

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Envs        []string
	Interactive bool
	Tty         bool
}

// NewCommand returns a new cobra.Command for exec
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	execCmd := &cobra.Command{
		Use:   "exec '[command-to-run --with-args]'",
		Short: "execute commands inside the uyuni containers using 'sh -c'",
		Run: func(cmd *cobra.Command, args []string) {
			run(globalFlags, flags, cmd, args)
		},
	}
	execCmd.Flags().StringArrayVarP(&flags.Envs, "env", "e", []string{}, "environment variables to pass to the command")
	execCmd.Flags().BoolVarP(&flags.Interactive, "interactive", "i", false, "Pass stdin to the container")
	execCmd.Flags().BoolVarP(&flags.Tty, "tty", "t", false, "Stdin is a TTY")
	return execCmd
}

func run(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {
	utils.Exec(globalFlags, flags.Interactive, flags.Tty, flags.Envs, args...)
}
