// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestRunner(t *testing.T) {
	type testCase struct {
		exit     int
		logLevel zerolog.Level
	}

	testCases := []testCase{
		{exit: 0, logLevel: zerolog.TraceLevel},
		{exit: 2, logLevel: zerolog.TraceLevel},
		{exit: 0, logLevel: zerolog.DebugLevel},
		{exit: 0, logLevel: zerolog.InfoLevel},
	}

	for i, test := range testCases {
		logWriter := new(strings.Builder)
		log.Logger = zerolog.New(logWriter)

		runner := NewRunner("sh", "-c",
			fmt.Sprintf(`echo "Test output: ENV=$ENV"; echo 'error message' >&2; exit %d`, test.exit),
		)
		out, err := runner.Log(test.logLevel).Env([]string{"ENV=foo"}).Exec()

		caseMsg := fmt.Sprintf("test %d: ", i)

		// Check the output
		testutils.AssertEquals(t, caseMsg+"Unexpected output", "Test output: ENV=foo\n", string(out))

		// Check the returned error
		if test.exit == 0 {
			testutils.AssertEquals(t, caseMsg+"Unexpected error", nil, err)
		} else {
			testutils.AssertEquals(t, caseMsg+"Unexpected error", "error message", string(err.Error()))
			var cmdErr *CmdError
			if errors.As(err, &cmdErr) {
				testutils.AssertEquals(t, caseMsg+"Unexpected exit code", test.exit, cmdErr.ExitCode())
			} else {
				t.Errorf(caseMsg + "unexpected error type")
			}
		}

		// Check the log content
		logContent := logWriter.String()
		t.Logf("log: %s", logContent)
		if test.logLevel == zerolog.TraceLevel {
			testutils.AssertTrue(t, caseMsg+"missing trace log entry", strings.Contains(logContent, "Command output:"))
		} else {
			testutils.AssertTrue(t, caseMsg+"unexpected trace log entry", !strings.Contains(logContent, `"level":"trace"`))
		}

		if test.logLevel <= zerolog.DebugLevel {
			testutils.AssertTrue(t, caseMsg+"missing debug log entry", strings.Contains(logContent, "Running:"))
		} else {
			testutils.AssertTrue(t, caseMsg+"unexpected debug log entry", !strings.Contains(logContent, `"level":"debug"`))
		}
	}
}

func ExampleRunner() {
	out, err := NewRunner("sh", "-c", `echo "Hello $user"`).
		Env([]string{"user=world"}).
		Log(zerolog.DebugLevel).
		Exec()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(strings.TrimSpace(string(out)))
	// Output: Hello world
}
