// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// OutputLogWriter contains information output the logger and the loglevel.
type OutputLogWriter struct {
	Logger   zerolog.Logger
	LogLevel zerolog.Level
}

// Write writes a byte array to an OutputLogWriter.
func (l OutputLogWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		// Trim CR added by stdlog.
		p = p[0 : n-1]
	}
	l.Logger.WithLevel(l.LogLevel).CallerSkipFrame(1).Msg(string(p))
	return
}

// RunCmd execute a shell command.
func RunCmd(command string, args ...string) error {
	log.Debug().Msgf("Running: %s %s", command, strings.Join(args, " "))

	return exec.Command(command, args...).Run()
}

// RunCmdStdMapping execute a shell command mapping the stdout and stderr.
func RunCmdStdMapping(logLevel zerolog.Level, command string, args ...string) error {
	localLogger := log.Level(logLevel)
	localLogger.Debug().Msgf("Running: %s %s", command, strings.Join(args, " "))

	runCmd := exec.Command(command, args...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	return runCmd.Run()
}

// RunCmdOutput execute a shell command and collects output.
func RunCmdOutput(logLevel zerolog.Level, command string, args ...string) ([]byte, error) {
	localLogger := log.Level(logLevel)

	localLogger.Debug().Msgf("Running: %s %s", command, strings.Join(args, " "))

	output, err := exec.Command(command, args...).Output()
	localLogger.Trace().Msgf("Command output: %s, error: %s", output, err)
	return output, err
}
