package uyunictl

import (
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/internal/utils"
)

var cpCmd = &cobra.Command{
	Use:   "cp [path/to/source.file] [path/to/desination.file]",
	Short: "copy files to and from the containers",
	Long: `copy takes a source and destination parameters.
One of them can be prefixed with 'server:' to indicate the path is within the server pod.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		command, podName := utils.GetPodName()
		commandArgs := []string{}
		extraArgs := []string{}
		src := strings.Replace(args[0], "server:", podName, 1)
		dst := strings.Replace(args[1], "server:", podName, 1)

		switch command {
		case "podman":
			commandArgs = []string{"cp", podName, src, dst}
		case "kubectl":
			commandArgs = []string{"cp", "-c", "uyuni", src, dst}
			extraArgs = []string{"-c", "uyuni", "--"}
		default:
			log.Fatalf("Unknown container kind: %s", command)
		}

		cp := exec.Command(command, commandArgs...)
		err := cp.Run()
		if err != nil {
			log.Fatalf("Failed to copy file: %s", err)
		}

		if user != "" && strings.HasPrefix(args[1], "server:") {
			execArgs := []string{"exec", podName}
			execArgs = append(execArgs, extraArgs...)
			owner := user
			if group != "" {
				owner = user + ":" + group
			}
			execArgs = append(execArgs, "chown", owner, src, dst)
			chown := exec.Command(command, execArgs...)
			err = chown.Run()
			if err != nil {
				log.Fatalf("Failed to change file owner: %s", err)
			}
		}
	},
}

var user string
var group string

func init() {
	cpCmd.Flags().StringVar(&user, "user", "", "User or UID to set on the destination file")
	cpCmd.Flags().StringVar(&user, "group", "", "Group or GID to set on the destination file")
	rootCmd.AddCommand(cpCmd)
}
