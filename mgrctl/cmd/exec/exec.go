// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Envs        []string `mapstructure:"env"`
	Interactive bool
	Tty         bool
	Backend     string
}

// NewCommand returns a new cobra.Command for exec.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags flagpole

	execCmd := &cobra.Command{
		Use:   "exec '[command-to-run --with-args]'",
		Short: L("Execute commands inside the uyuni containers using 'sh -c'"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, run)
		},
	}

	execCmd.Flags().StringSliceP("env", "e", []string{}, L("environment variables to pass to the command, separated by commas"))
	execCmd.Flags().BoolP("interactive", "i", false, L("Pass stdin to the container"))
	execCmd.Flags().BoolP("tty", "t", false, L("Stdin is a TTY"))

	utils.AddBackendFlag(execCmd)
	return execCmd
}

func run(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)
	podName, err := cnx.GetPodName()
	if err != nil {
		log.Fatal().Err(err)
	}

	command, err := cnx.GetCommand()
	if err != nil {
		log.Fatal().Err(err)
	}

	commandArgs := []string{"exec"}
	envs := []string{}
	envs = append(envs, flags.Envs...)
	if flags.Interactive {
		commandArgs = append(commandArgs, "-i")
		envs = append(envs, "ENV=/etc/sh.shrc.local")
	}
	if flags.Tty {
		commandArgs = append(commandArgs, "-t")
		envs = append(envs, "TERM")
	}
	commandArgs = append(commandArgs, podName)

	if command == "kubectl" {
		commandArgs = append(commandArgs, "-c", "uyuni", "--")
	}

	newEnv := []string{}
	for _, envValue := range envs {
		if !strings.Contains(envValue, "=") {
			if value, set := os.LookupEnv(envValue); set {
				newEnv = append(newEnv, fmt.Sprintf("%s=%s", envValue, value))
			}
		} else {
			newEnv = append(newEnv, envValue)
		}
	}
	if len(newEnv) > 0 {
		commandArgs = append(commandArgs, "env")
		commandArgs = append(commandArgs, newEnv...)
	}
	commandArgs = append(commandArgs, "sh", "-c", strings.Join(args, " "))
	err = RunRawCmd(command, commandArgs)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Info().Err(err).Msg(L("Command failed"))
			os.Exit(exitErr.ExitCode())
		}
	}
	log.Info().Msg(L("Command returned with exit code 0"))

	return nil
}

type copyWriter struct {
	Stream io.Writer
}

// Write writes an array of buffer in a stream.
func (l copyWriter) Write(p []byte) (n int, err error) {
	// Filter out kubectl line about terminated exit code
	if !strings.HasPrefix(string(p), "command terminated with exit code") {
		if _, err := l.Stream.Write(p); err != nil {
			return 0, fmt.Errorf(L("cannot write: %s"), err)
		}

		n = len(p)
		if n > 0 && p[n-1] == '\n' {
			// Trim CR added by stdlog.
			p = p[0 : n-1]
		}
		log.Debug().Msg(string(p))
	}
	return
}

// RunRawCmd runs a command, mapping stdout and start error, waiting and checking return code.
func RunRawCmd(command string, args []string) error {
	log.Info().Msgf(L("Running: %s %s"), command, strings.Join(args, " "))

	runCmd := exec.Command(command, args...)
	runCmd.Stdin = os.Stdin

	runCmd.Stdout = copyWriter{Stream: os.Stdout}
	runCmd.Stderr = copyWriter{Stream: os.Stderr}

	if err := runCmd.Start(); err != nil {
		log.Debug().Err(err).Msg("error starting command")
		return err
	}

	return runCmd.Wait()
}
