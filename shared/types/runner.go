// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"bytes"

	"github.com/rs/zerolog"
)

// Runner is an interface to execute system calls.
type Runner interface {
	// Log sets the log level of the output.
	Log(logLevel zerolog.Level) Runner

	// Spinner sets a spinner with its message.
	// If no message is passed, the command will be used.
	Spinner(message string) Runner

	// StdMapping maps the process output and error streams to the standard ones.
	// This is useful to show the process output in the console and the logs and can be combined with Log().
	StdMapping() Runner

	// Std maps the process output to the out bytes buffer.
	// This is useful to get process output for backgrounds tasks.
	Std(out *bytes.Buffer) Runner

	// Wait waits for the command running in the background to ends.
	Wait() error

	// InputString adds a string as input of the process.
	InputString(input string) Runner

	// Env sets environment variables to use for the command.
	Env(env []string) Runner

	// Exec really executes the command and returns its output and error.
	// The error output to used as error message if the StdMapping() function wasn't called.
	Exec() ([]byte, error)

	// Start starts the command, particularly for commands to run in the background.
	Start() error
}
