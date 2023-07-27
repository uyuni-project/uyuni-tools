package utils

import (
	"log"
	"strings"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// Copy transfers a file to or from the container.
// Prefix one of src or dst parameters with `server:` to designate the path is in the container
// user and group parameters are used to set the owner of a file transfered in the container.
func Copy(globalFlags *types.GlobalFlags, src string, dst string, user string, group string) {
	command, podName := GetPodName()
	var commandArgs []string
	extraArgs := []string{}
	srcExpanded := strings.Replace(src, "server:", podName+":", 1)
	dstExpanded := strings.Replace(dst, "server:", podName+":", 1)

	switch command {
	case "podman":
		commandArgs = []string{"cp", srcExpanded, dstExpanded}
	case "kubectl":
		commandArgs = []string{"cp", "-c", "uyuni", srcExpanded, dstExpanded}
		extraArgs = []string{"-c", "uyuni", "--"}
	default:
		log.Fatalf("Unknown container kind: %s\n", command)
	}

	RunCmd(command, commandArgs, "Failed to copy file", globalFlags.Verbose)

	if user != "" && strings.HasPrefix(dst, "server:") {
		execArgs := []string{"exec", podName}
		execArgs = append(execArgs, extraArgs...)
		owner := user
		if group != "" {
			owner = user + ":" + group
		}
		execArgs = append(execArgs, "chown", owner, strings.Replace(dst, "server:", "", 1))
		RunCmd(command, execArgs, "Failed to change file owner", globalFlags.Verbose)
	}
}
