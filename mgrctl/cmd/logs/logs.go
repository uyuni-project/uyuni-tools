// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package logs

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/exec"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Taskomatic bool
	Web        bool
	Salt       bool
	Reposync   bool
	Files      string
	Follow     bool
	Backend    string
}

// NewCommand returns a new cobra.Command for logs.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags flagpole

	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: L("Show or follow logs of Uyuni services inside the container"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	logsCmd.Flags().Bool("taskomatic", false, L("Show Taskomatic logs"))
	logsCmd.Flags().Bool("web", false, L("Show Web UI logs"))
	logsCmd.Flags().Bool("salt", false, L("Show Salt logs"))
	logsCmd.Flags().Bool("reposync", false, L("Show Reposync logs"))
	logsCmd.Flags().String("files", "", L("Regex to filter reposync log files"))
	logsCmd.Flags().BoolP("follow", "f", false, L("Follow log output"))

	utils.AddBackendFlag(logsCmd)
	return logsCmd
}

func run(_ *types.GlobalFlags, flags *flagpole, _ *cobra.Command, _ []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)
	podName, err := cnx.GetPodName()
	if err != nil {
		return err
	}

	command, err := cnx.GetCommand()
	if err != nil {
		return err
	}

	var paths []string
	if flags.Taskomatic {
		paths = append(paths, "/var/log/rhn/rhn_tasko*.log")
	}
	if flags.Web {
		paths = append(paths, "/var/log/rhn/rhn_web*.log")
	}
	if flags.Salt {
		paths = append(paths, "/var/log/salt/api", "/var/log/salt/master", "/var/log/salt/minion")
	}
	if flags.Reposync {
		if flags.Files != "" {
			findCmd := fmt.Sprintf("find /var/log/rhn/reposync/ -type f | grep -E '%s'", flags.Files)
			paths = append(paths, fmt.Sprintf("$(%s)", findCmd))
		} else {
			paths = append(paths, "/var/log/rhn/reposync/*.log")
		}
	}

	if len(paths) == 0 {
		return fmt.Errorf(L("please specify at least one service to get logs for (e.g., --web, --salt)"))
	}

	tailCmd := "tail"
	if flags.Follow {
		tailCmd += " -f"
	}

	shCmd := fmt.Sprintf("%s %s", tailCmd, strings.Join(paths, " "))

	commandArgs := []string{"exec", "-i", "-t"}

	if command == "kubectl" {
		namespace, err := cnx.GetNamespace("")
		if err != nil {
			return err
		}
		commandArgs = append(commandArgs, "-n", namespace, "-c", "uyuni", podName, "--")
	} else {
		commandArgs = append(commandArgs, podName)
	}

	commandArgs = append(commandArgs, "bash", "-c", shCmd)

	// Re-use the execution logic from the exec package
	return exec.RunRawCmd(command, commandArgs)
}
