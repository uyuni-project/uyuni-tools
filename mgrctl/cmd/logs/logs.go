// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package logs

import (
	"errors"
	"fmt"
	"os"
	os_exec "os/exec"
	"strings"

	"github.com/spf13/cobra"
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

// runCmd allows unit tests to mock or skip the main execution loop.
var runCmd = run

// NewCommand returns a new cobra.Command for logs.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags flagpole

	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: L("Show or follow logs of Uyuni services inside the container"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, runCmd)
		},
	}

	logsCmd.Flags().Bool("taskomatic", false, L("Show Taskomatic logs"))
	logsCmd.Flags().Bool("web", false, L("Show Web UI logs"))
	logsCmd.Flags().Bool("salt", false, L("Show Salt logs"))
	logsCmd.Flags().Bool("reposync", false, L("Show Reposync logs"))
	logsCmd.Flags().String("files", "",
		L("Regular expression to filter reposync log files (matches against the full file path)"))
	logsCmd.Flags().BoolP("follow", "f", false, L("Follow log output"))

	utils.AddBackendFlag(logsCmd)
	return logsCmd
}

func validateFlags(flags *flagpole) error {
	if flags.Files != "" && !flags.Reposync {
		return errors.New(L("--files flag can only be used with --reposync"))
	}
	return nil
}

func getLogPaths(flags *flagpole) []string {
	var paths []string
	if flags.Taskomatic {
		paths = append(paths, "/var/log/rhn/rhn_tasko*.log")
	}
	if flags.Web {
		paths = append(paths, "/var/log/rhn/rhn_web*.log")
	}
	if flags.Salt {
		// Minion does not run inside the Uyuni server container
		paths = append(paths, "/var/log/salt/api", "/var/log/salt/master")
	}
	if flags.Reposync {
		if flags.Files != "" {
			script := "files=$(find /var/log/rhn/reposync/ -type f | grep -E %q); " +
				"if [ -z \"$files\" ]; then "
			// Use literal path on follow so tail can wait for file creation
			if flags.Follow {
				script += "echo \"/var/log/rhn/reposync/*.log\"; else echo $files; fi"
			} else {
				script += "echo 'No matching reposync logs found.' >&2; exit 1; fi; echo $files"
			}
			paths = append(paths, fmt.Sprintf("$(%s)", fmt.Sprintf(script, flags.Files)))
		} else {
			paths = append(paths, "/var/log/rhn/reposync/*.log")
		}
	}
	return paths
}

func run(_ *types.GlobalFlags, flags *flagpole, _ *cobra.Command, _ []string) error {
	if err := validateFlags(flags); err != nil {
		return err
	}

	paths := getLogPaths(flags)
	if len(paths) == 0 {
		return errors.New(L("please specify at least one service to get logs for (e.g., --web, --salt)"))
	}

	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)
	podName, err := cnx.GetPodName()
	if err != nil {
		return err
	}

	command, err := cnx.GetCommand()
	if err != nil {
		return err
	}

	tailCmd := "tail"
	if flags.Follow {
		tailCmd += " -F"
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

	_, err = utils.NewRunner(command, commandArgs...).Exec()
	if err != nil {
		var exitErr *os_exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		return err
	}
	return nil
}
