// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"crypto/rand"
	"encoding/hex"
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

func prepareSource(args []string, cnx *shared.Connection) (string, error) {
	source := "-"
	if len(args) > 0 {
		source = args[0]
		if !utils.FileExists(source) {
			return "", fmt.Errorf(L("source %s does not exists"), source)
		}
		randBytes := make([]byte, 16)
		if _, err := rand.Read(randBytes); err != nil {
			return "", fmt.Errorf(L("unable to get random file prefix: %s"), err)
		}
		source = hex.EncodeToString(randBytes) + source
		if err := cnx.Copy(args[0], "server:"+source, "", ""); err != nil {
			return "", err
		}
	}
	return source, nil
}

func cleanupSource(file string, cnx *shared.Connection) {
	if _, err := cnx.Exec("rm", file); err != nil {
		log.Error().Err(err).Msg(L("unable to cleanup source file"))
	}
}

func prepareOutput(flags *configFlags) (string, error) {
	output := "-"
	if flags.OutputFile != "" {
		output = flags.OutputFile
		if utils.FileExists(output) && !flags.ForceOverwrite {
			return "", fmt.Errorf(L("output file %s exists, use -f to force overwrite"), output)
		}
	}
	return output, nil
}

func getBaseCommand(keepStdin bool, flags *configFlags, cnx *shared.Connection) (string, []string, error) {
	podName, err := cnx.GetPodName()
	if err != nil {
		return "", nil, err
	}

	command, err := cnx.GetCommand()
	if err != nil {
		return "", nil, err
	}

	commandArgs := []string{"exec"}
	envs := []string{}
	if flags.Interactive {
		commandArgs = append(commandArgs, "-i")
		envs = append(envs, "ENV=/etc/sh.shrc.local")
		commandArgs = append(commandArgs, "-t")
		envs = append(envs, "TERM")
	} else if keepStdin {
		// To use STDIN source, we need to pass -i
		commandArgs = append(commandArgs, "-i")
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
	return command, commandArgs, nil
}

func doSql(globalFlags *types.GlobalFlags, flags *configFlags, cmd *cobra.Command, args []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)

	// Validate options
	source, err := prepareSource(args, cnx)
	if err != nil {
		return err
	}
	if source != "" && source != "-" {
		defer cleanupSource(source, cnx)
	}
	output, err := prepareOutput(flags)
	if err != nil {
		return err
	}

	// For now do quick wrapper around spacewalk-sql tool.
	// TODO - ideally use sql directly, but will need some gateway to be able to connect to the database
	command, commandArgs, err := getBaseCommand(source == "-", flags, cnx)
	if err != nil {
		return err
	}
	commandArgs = append(commandArgs, "/usr/bin/spacewalk-sql")

	sqlArgs := []string{}
	if flags.Database == "reportdb" {
		sqlArgs = append(sqlArgs, "--reportdb")
	} else if flags.Database != "productdb" {
		return fmt.Errorf(L("unknown or unsupported database %s"), flags.Database)
	}

	if flags.Interactive {
		sqlArgs = append(sqlArgs, "-i")
	} else {
		sqlArgs = append(sqlArgs, "--select-mode", source)
	}
	commandArgs = append(commandArgs, sqlArgs...)

	err = runCmd(command, output, commandArgs)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Info().Err(err).Msg(L("Command failed"))
			os.Exit(exitErr.ExitCode())
		}
	}
	if output != "-" {
		log.Info().Msgf(L("Result is stored in the file '%s'"), output)
	}
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

func runCmd(command string, output string, args []string) error {
	log.Info().Msgf(L("Running: %s %s"), command, strings.Join(args, " "))

	runCmd := exec.Command(command, args...)
	runCmd.Stdin = os.Stdin

	if output == "" || output == "-" {
		runCmd.Stdout = copyWriter{Stream: os.Stdout}
	} else {
		log.Trace().Msgf("Output is FILE %s", output)
		f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		runCmd.Stdout = copyWriter{Stream: f}
	}
	runCmd.Stderr = copyWriter{Stream: os.Stderr}

	if err := runCmd.Start(); err != nil {
		log.Debug().Err(err).Msg("error starting command")
		return err
	}

	return runCmd.Wait()
}
