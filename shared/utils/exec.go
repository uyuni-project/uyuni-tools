package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func Exec(globalFlags *types.GlobalFlags, backend string, interactive bool, tty bool, outputToLog bool, env []string, args ...string) {
	command, podName := GetPodName(globalFlags, backend, true)

	commandArgs := []string{"exec"}
	if interactive {
		commandArgs = append(commandArgs, "-i")
	}
	if tty {
		commandArgs = append(commandArgs, "-t")
	}
	commandArgs = append(commandArgs, podName)

	if command == "kubectl" {
		commandArgs = append(commandArgs, "-c", "uyuni", "--")
	}

	newEnv := []string{}
	for _, envValue := range env {
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

	err := RunRawCmd(command, commandArgs, outputToLog)
	if err != nil {
		log.Fatal().Err(err).Msg("error running the command")
	}
}

// should we call fatal or let the caller cal it?
func RunRawCmd(command string, args []string, outputToLog bool) error {

	// TODO think if we should log all command or just the sub-command part
	// most problematic part can be the flags
	log.Debug().Msgf(" Running: %s %s", command, strings.Join(args, " "))

	runCmd := exec.Command(command, args...)
	runCmd.Stdin = os.Stdin
	if outputToLog {
		runCmd.Stdout = log.Logger
	} else {
		runCmd.Stdout = os.Stdout
	}

	// Filter out kubectl line about terminated exit code
	stderr, err := runCmd.StderrPipe()
	if err != nil {
		log.Debug().Err(err).Msg("error starting stderr processor for command")
		return err
	}
	defer stderr.Close()

	if err = runCmd.Start(); err != nil {
		log.Debug().Err(err).Msg("error starting command")
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "command terminated with exit code") {
			if outputToLog {
				log.Error().Msg(line)
			} else {
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}

	if scanner.Err() != nil {
		log.Debug().Msg("error scanning stderr")
		return scanner.Err()
	}
	if err = runCmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok {
			log.Debug().Err(exitErr).Msgf("error on command exit code")
			return exitErr
		}
		log.Debug().Err(err).Msg("error on wait command")
		return err
	}
	return nil
}
