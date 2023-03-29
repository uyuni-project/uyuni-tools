package cp

import (
	"log"
	"strings"

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
		Long: `copy takes a source and destination parameters.
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
	command, podName := utils.GetPodName()
	commandArgs := []string{}
	extraArgs := []string{}
	src := strings.Replace(args[0], "server:", podName+":", 1)
	dst := strings.Replace(args[1], "server:", podName+":", 1)

	switch command {
	case "podman":
		commandArgs = []string{"cp", podName, src, dst}
	case "kubectl":
		commandArgs = []string{"cp", "-c", "uyuni", src, dst}
		extraArgs = []string{"-c", "uyuni", "--"}
	default:
		log.Fatalf("Unknown container kind: %s\n", command)
	}

	utils.RunCmd(command, commandArgs, "Failed to copy file", globalFlags.Verbose)

	if flags.User != "" && strings.HasPrefix(args[1], "server:") {
		execArgs := []string{"exec", podName}
		execArgs = append(execArgs, extraArgs...)
		owner := flags.User
		if flags.Group != "" {
			owner = flags.User + ":" + flags.Group
		}
		execArgs = append(execArgs, "chown", owner, strings.Replace(args[1], "server:", "", 1))
		utils.RunCmd(command, execArgs, "Failed to change file owner", globalFlags.Verbose)
	}
}
