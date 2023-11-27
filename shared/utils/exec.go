// SPDX-FileCopyrightText: 2023 SUSE LLC
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

type OutputLogWriter struct {
	Logger   zerolog.Logger
	LogLevel zerolog.Level
}

func (l OutputLogWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		// Trim CR added by stdlog.
		p = p[0 : n-1]
	}
	l.Logger.WithLevel(l.LogLevel).CallerSkipFrame(1).Msg(string(p))
	return
}

func RunCmd(command string, args ...string) error {
	log.Debug().Msgf("Running: %s %s", command, strings.Join(args, " "))

	return exec.Command(command, args...).Run()
}

func RunCmdStdMapping(command string, args ...string) error {
	log.Debug().Msgf("Running: %s %s", command, strings.Join(args, " "))

	runCmd := exec.Command(command, args...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	return runCmd.Run()
}

func RunCmdOutput(logLevel zerolog.Level, command string, args ...string) ([]byte, error) {
	logger := log.Logger.WithLevel(logLevel)

	logger.Msgf("Running: %s %s", command, strings.Join(args, " "))

	output, err := exec.Command(command, args...).Output()
	if err != nil {
		logger.Err(err).Msgf("Command returned Error: %s", output)
	} else if len(output) > 0 {
		logger.Msgf("Command output: %s", output)
	}
	return output, err
}
