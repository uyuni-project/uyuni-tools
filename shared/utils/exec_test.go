// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestRunner(t *testing.T) {
	type testCase struct {
		name     string
		exit     int
		logLevel zerolog.Level
	}

	testCases := []testCase{
		{name: "exit 0 with trace log level", exit: 0, logLevel: zerolog.TraceLevel},
		{name: "exit 2 with trace log level", exit: 2, logLevel: zerolog.TraceLevel},
		{name: "exit 0 with debug log level", exit: 0, logLevel: zerolog.DebugLevel},
		{name: "exit 0 with info log level", exit: 0, logLevel: zerolog.InfoLevel},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			logWriter := new(strings.Builder)
			log.Logger = zerolog.New(logWriter)

			runner := NewRunner("sh", "-c",
				fmt.Sprintf(`echo "Test output: ENV=$ENV"; echo 'error message' >&2; exit %d`, test.exit),
			)
			out, err := runner.Log(test.logLevel).Env([]string{"ENV=foo"}).Exec()

			// Check the output
			testutils.AssertEquals(t, "Unexpected output", "Test output: ENV=foo\n", string(out))

			// Check the returned error
			if test.exit == 0 {
				testutils.AssertNoError(t, "Unexpected error", err)
			} else {
				testutils.AssertError(t, "error message", err)
				var cmdErr *CmdError
				if errors.As(err, &cmdErr) {
					testutils.AssertEquals(t, "Unexpected exit code", test.exit, cmdErr.ExitCode())
				} else {
					t.Errorf("unexpected error type")
				}
			}

			// Check the log content
			logContent := logWriter.String()
			t.Logf("log: %s", logContent)
			if test.logLevel == zerolog.TraceLevel {
				testutils.AssertStringContains(t, "missing trace log entry", logContent, "Command output:")
			} else {
				testutils.AssertTrue(t, "unexpected trace log entry", !strings.Contains(logContent, `"level":"trace"`))
			}

			if test.logLevel <= zerolog.DebugLevel {
				testutils.AssertStringContains(t, "missing debug log entry", logContent, "Running:")
			} else {
				testutils.AssertTrue(t, "unexpected debug log entry", !strings.Contains(logContent, `"level":"debug"`))
			}
		})
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

func TestContextRunner(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	testData := []byte("Hello world")

	runner := NewRunnerWithContext(ctx, "sh", "-c", fmt.Sprintf(`echo -n "%s"`, testData))
	var out bytes.Buffer
	err := runner.Std(&out).Start()
	testutils.AssertNoError(t, "Unexpected start failure", err)
	err = runner.Wait()
	testutils.AssertNoError(t, "Unexpected wait failure", err)

	testutils.AssertEquals(t, "Output does not match", out.Bytes(), testData)
}
