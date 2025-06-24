// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/types"
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

// NewRunner creates a new runner instance for the command.
func NewRunner(command string, args ...string) types.Runner {
	runner := runnerImpl{logger: log.Logger}
	runner.cmd = exec.Command(command, args...)
	return &runner
}

// runnerImpl is a helper object around the exec.Command() function.
// It implements the Runner interface.
//
// This is supposed to be created using the NewRunner() function.
type runnerImpl struct {
	logger  zerolog.Logger
	cmd     *exec.Cmd
	spinner *spinner.Spinner
}

// Log sets the log level of the output.
func (r *runnerImpl) Log(logLevel zerolog.Level) types.Runner {
	r.logger = log.Logger.Level(logLevel)
	return r
}

// Spinner sets a spinner with its message.
// If no message is passed, the command will be used.
func (r *runnerImpl) Spinner(message string) types.Runner {
	r.spinner = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	text := message
	if message == "" {
		text = strings.Join(r.cmd.Args, " ")
	}
	r.spinner.Suffix = fmt.Sprintf(" %s\n", text)
	return r
}

// StdMapping maps the process output and error streams to the standard ones.
// This is useful to show the process output in the console and the logs and can be combined with Log().
func (r *runnerImpl) StdMapping() types.Runner {
	r.cmd.Stdout = r.logger
	r.cmd.Stderr = r.logger
	return r
}

// Env sets environment variables to use for the command.
func (r *runnerImpl) Env(env []string) types.Runner {
	if r.cmd.Env == nil {
		r.cmd.Env = os.Environ()
	}
	r.cmd.Env = append(r.cmd.Env, env...)
	return r
}

// Exec really executes the command and returns its output and error.
// The error output to used as error message if the StdMapping() function wasn't called.
func (r *runnerImpl) Exec() ([]byte, error) {
	if r.spinner != nil {
		r.spinner.Start()
	}

	r.logger.Debug().Msgf("Running: %s", strings.Join(r.cmd.Args, " "))
	var out []byte
	var err error

	if r.cmd.Stdout != nil {
		err = r.cmd.Run()
	} else {
		out, err = r.cmd.Output()
	}

	if r.spinner != nil {
		r.spinner.Stop()
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		err = &CmdError{exitErr}
	}

	r.logger.Trace().Msgf("Command output: %s, error: %s", out, err)

	return out, err
}

// CmdError is a wrapper around exec.ExitError to show the standard error as message.
type CmdError struct {
	*exec.ExitError
}

// Error returns the stderr as error message.
func (e *CmdError) Error() string {
	return strings.TrimSpace(string(e.Stderr))
}

// RunCmd execute a shell command.
func RunCmd(command string, args ...string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Suffix = fmt.Sprintf(" %s %s\n", command, strings.Join(args, " "))
	s.Start() // Start the spinner
	log.Debug().Msgf("Running: %s %s", command, strings.Join(args, " "))
	err := exec.Command(command, args...).Run()
	s.Stop()
	return err
}

// RunCmdStdMapping execute a shell command mapping the stdout and stderr.
func RunCmdStdMapping(logLevel zerolog.Level, command string, args ...string) error {
	localLogger := log.Logger.Level(logLevel)
	localLogger.Debug().Msgf("Running: %s %s", command, strings.Join(args, " "))

	runCmd := exec.Command(command, args...)
	runCmd.Stdout = localLogger
	runCmd.Stderr = localLogger
	err := runCmd.Run()
	return err
}

// RunCmdOutput execute a shell command and collects output.
func RunCmdOutput(logLevel zerolog.Level, command string, args ...string) ([]byte, error) {
	localLogger := log.Level(logLevel)
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Suffix = fmt.Sprintf(" %s %s\n", command, strings.Join(args, " "))
	if logLevel != zerolog.Disabled {
		s.Start() // Start the spinner
	}
	localLogger.Debug().Msgf("Running: %s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	output, err := cmd.Output()
	if logLevel != zerolog.Disabled {
		s.Stop()
	}
	localLogger.Trace().Msgf("Command output: %s, error: %s", output, err)
	message := strings.TrimSpace(errBuf.String())
	if message != "" {
		err = errors.New(message)
	}
	return output, err
}

// RunCmdInput execute a shell command and pass input string to the StdIn.
func RunCmdInput(command string, input string, args ...string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Suffix = fmt.Sprintf(" %s %s\n", command, strings.Join(args, " "))
	s.Start() // Start the spinner
	log.Debug().Msgf("Running: %s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Stdin = strings.NewReader(input)
	err := cmd.Run()
	s.Stop()
	return err
}

// IsInstalled checks if a tool is in the path.
func IsInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// GetEnvironmentVarsList returns list of environmental variables to be passed to exec.
func GetEnvironmentVarsList() []string {
	// Taken from /etc/profile and /etc/profile.d/lang
	return []string{"TERM", "PAGER",
		"LESS", "LESSOPEN", "LESSKEY", "LESSCLOSE", "LESS_ADVANCED_PREPROCESSOR", "MORE",
		"LANG", "LC_CTYPE", "LC_ALL"}
}
