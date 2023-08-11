package cp

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	User  string
	Group string
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	cpCmd := &cobra.Command{
		Use:   "cp [path/to/source.file] [path/to/desination.file]",
		Short: "copy files to and from the containers",
		Long: `Takes a source and destination parameters.
	One of them can be prefixed with 'server:' to indicate the path is within the server pod.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			run(globalFlags, flags, cmd, args)
		},
	}

	cpCmd.Flags().StringVar(&flags.User, "user", "", "User or UID to set on the destination file")
	cpCmd.Flags().StringVar(&flags.Group, "group", "", "Group or GID to set on the destination file")
	return cpCmd
}

func run(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {
	utils.Copy(globalFlags, args[0], args[1], flags.User, flags.Group)
}
