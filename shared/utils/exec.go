package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

type outputLogWriter struct {
	logger zerolog.Logger
}

func (l outputLogWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		// Trim CR added by stdlog.
		p = p[0 : n-1]
	}
	l.logger.Debug().CallerSkipFrame(1).Msg(string(p))
	return
}

func Exec(globalFlags *types.GlobalFlags, backend string, interactive bool, tty bool, 
	outputToLog bool, env []string, args ...string) error{
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

	return RunRawCmd(command, commandArgs, outputToLog)
	
}

func RunRawCmd(command string, args []string, outputToLog bool) error {

	log.Debug().Msgf(" Running: %s %s", command, strings.Join(args, " "))

	runCmd := exec.Command(command, args...)
	runCmd.Stdin = os.Stdin
	if outputToLog {
		runCmd.Stdout = outputLogWriter{log.Logger}
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
		// needed because of kubernetes installation, to ignore the stderr
		if !strings.HasPrefix(line, "command terminated with exit code") {
			if outputToLog {
				log.Debug().Msg(line)
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

func RunCmdOutput(command string, args ...string) ([]byte, error) {

	log.Debug().Msgf(" Running: %s %s", command, strings.Join(args, " "))

	output, err := exec.Command(command, args...).Output()
	if err != nil {
		log.Debug().Err(err).Msgf("Command returned Error: %s", output)
	} else if len(output) > 0 {
		log.Debug().Msgf("Command output: %s", output)
	}
	return output, err
}
