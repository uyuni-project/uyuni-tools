// SPDX-FileCopyrightText: 2026 Shubham Kapoor <165891311+shubhamk0205@users.noreply.github.com>
//
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"errors"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestRunReportsStartError(t *testing.T) {
	originalRunRawCmd := runRawCmd
	defer func() { runRawCmd = originalRunRawCmd }()

	runRawCmd = func(_ string, _ []string) error {
		return errors.New("executable file not found in $PATH")
	}

	err := run(nil, &flagpole{Backend: "host"}, nil, []string{"echo", "hello"})
	testutils.AssertError(t, "cannot execute host: executable file not found in $PATH", err)
}

func TestRunSuccess(t *testing.T) {
	originalRunRawCmd := runRawCmd
	defer func() { runRawCmd = originalRunRawCmd }()

	var calledCommand string
	var calledArgs []string
	runRawCmd = func(command string, args []string) error {
		calledCommand = command
		calledArgs = args
		return nil
	}

	err := run(nil, &flagpole{Backend: "host"}, nil, []string{"echo", "hello"})
	testutils.AssertNoError(t, "unexpected error", err)
	testutils.AssertEquals(t, "wrong command", "host", calledCommand)
	testutils.AssertEquals(t, "wrong arguments", "exec host sh -c echo hello", strings.Join(calledArgs, " "))
}
