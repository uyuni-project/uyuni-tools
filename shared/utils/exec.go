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

func Exec(globalFlags *types.GlobalFlags, backend string, interactive bool, tty bool, env []string, args ...string) {
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
	if globalFlags.Verbose {
		log.Info().Msgf(" Running: %s %s", command, strings.Join(commandArgs, " "))
	}
	runCmd := exec.Command(command, commandArgs...)
	runCmd.Stdout = os.Stdout
	runCmd.Stdin = os.Stdin

	// Filter out kubectl line about terminated exit code
	stderr, err := runCmd.StderrPipe()
	if err != nil {
		log.Fatal().Err(err).Msg("error running command")
	}
	if err = runCmd.Start(); err != nil {
		log.Fatal().Err(err).Msg("error running command")
	}
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "command terminated with exit code") {
			fmt.Fprintln(os.Stderr, line)
		}
	}

	if scanner.Err() != nil {
		log.Fatal().Err(scanner.Err()).Msg("error running command")
	}
	if err = runCmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatal().Err(err).Msg("error running command")
		}
	}
}
