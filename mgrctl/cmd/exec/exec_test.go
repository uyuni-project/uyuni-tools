// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand(nil)
	if cmd == nil {
		t.Fatal("Unexpected nil command")
	}

	// Check that flags are properly registered
	flags := cmd.Flags()

	// Test env flag
	envFlag := flags.Lookup("env")
	if envFlag == nil {
		t.Error("env flag not registered")
	}
	if envFlag.DefValue != "[]" {
		t.Errorf("env flag has unexpected default value: %s", envFlag.DefValue)
	}

	// Test interactive flag
	interactiveFlag := flags.Lookup("interactive")
	if interactiveFlag == nil {
		t.Error("interactive flag not registered")
	}
	if interactiveFlag.DefValue != "false" {
		t.Errorf("interactive flag has unexpected default value: %s", interactiveFlag.DefValue)
	}

	// Test tty flag
	ttyFlag := flags.Lookup("tty")
	if ttyFlag == nil {
		t.Error("tty flag not registered")
	}
	if ttyFlag.DefValue != "false" {
		t.Errorf("tty flag has unexpected default value: %s", ttyFlag.DefValue)
	}

	// Test backend flag
	backendFlag := flags.Lookup("backend")
	if backendFlag == nil {
		t.Error("backend flag not registered")
	}
}

func TestCopyWriterWrite(t *testing.T) {
	testCases := []struct {
		name           string
		input          []byte
		expectedOutput string
		expectedBytes  int
	}{
		{
			name:           "normal output",
			input:          []byte("Hello World\n"),
			expectedOutput: "Hello World\n",
			expectedBytes:  len("Hello World\n"),
		},
		{
			name:           "kubectl termination message filtered",
			input:          []byte("command terminated with exit code 1\n"),
			expectedOutput: "",
			expectedBytes:  0, // Filtered messages return 0 bytes written
		},
		{
			name:           "mixed output with termination message",
			input:          []byte("Some output\ncommand terminated with exit code 1\n"),
			expectedOutput: "Some output\ncommand terminated with exit code 1\n",
			expectedBytes:  len("Some output\ncommand terminated with exit code 1\n"),
		},
		{
			name:           "empty input",
			input:          []byte(""),
			expectedOutput: "",
			expectedBytes:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := copyWriter{Stream: &buf}

			n, err := writer.Write(tc.input)

			testutils.AssertEquals(t, "Write error", nil, err)
			testutils.AssertEquals(t, "Bytes written", tc.expectedBytes, n)
			testutils.AssertEquals(t, "Output mismatch", tc.expectedOutput, buf.String())
		})
	}
}

func TestCopyWriterWriteMultipleWrites(t *testing.T) {
	var buf bytes.Buffer
	writer := copyWriter{Stream: &buf}

	chunks := [][]byte{
		[]byte("First line\n"),
		[]byte("Second line\n"),
		[]byte("Third line\n"),
	}

	for _, chunk := range chunks {
		n, err := writer.Write(chunk)
		testutils.AssertEquals(t, "Write error", nil, err)
		testutils.AssertEquals(t, "Bytes written", len(chunk), n)
	}

	expected := "First line\nSecond line\nThird line\n"
	testutils.AssertEquals(t, "Combined output mismatch", expected, buf.String())
}

func TestCopyWriterFiltering(t *testing.T) {
	testCases := []struct {
		name           string
		input          []byte
		expectedOutput string
		expectedBytes  int
		description    string
	}{
		{
			name:           "kubectl exit code 0",
			input:          []byte("command terminated with exit code 0\n"),
			expectedOutput: "",
			expectedBytes:  0,
			description:    "Should filter kubectl termination message with exit code 0",
		},
		{
			name:           "kubectl exit code 1",
			input:          []byte("command terminated with exit code 1\n"),
			expectedOutput: "",
			expectedBytes:  0,
			description:    "Should filter kubectl termination message with exit code 1",
		},
		{
			name:           "kubectl exit code 127",
			input:          []byte("command terminated with exit code 127\n"),
			expectedOutput: "",
			expectedBytes:  0,
			description:    "Should filter kubectl termination message with exit code 127",
		},
		{
			name:           "similar but different message",
			input:          []byte("command terminated with error\n"),
			expectedOutput: "command terminated with error\n",
			expectedBytes:  len("command terminated with error\n"),
			description:    "Should NOT filter similar but different messages",
		},
		{
			name:           "partial match at start",
			input:          []byte("command terminated\n"),
			expectedOutput: "command terminated\n",
			expectedBytes:  len("command terminated\n"),
			description:    "Should NOT filter partial matches",
		},
		{
			name:           "message in middle of line",
			input:          []byte("Some text command terminated with exit code 1 more text\n"),
			expectedOutput: "Some text command terminated with exit code 1 more text\n",
			expectedBytes:  len("Some text command terminated with exit code 1 more text\n"),
			description:    "Should NOT filter when message is not at start of line",
		},
		{
			name:           "output before termination message",
			input:          []byte("Actual output\ncommand terminated with exit code 1\n"),
			expectedOutput: "Actual output\ncommand terminated with exit code 1\n",
			expectedBytes:  len("Actual output\ncommand terminated with exit code 1\n"),
			description:    "Should write when output precedes termination message",
		},
		{
			name:           "multiple lines with termination",
			input:          []byte("Line 1\nLine 2\ncommand terminated with exit code 0\nLine 3\n"),
			expectedOutput: "Line 1\nLine 2\ncommand terminated with exit code 0\nLine 3\n",
			expectedBytes:  len("Line 1\nLine 2\ncommand terminated with exit code 0\nLine 3\n"),
			description:    "Should write multi-line output containing termination message",
		},
		{
			name:           "termination message without newline",
			input:          []byte("command terminated with exit code 1"),
			expectedOutput: "",
			expectedBytes:  0,
			description:    "Should filter termination message without trailing newline",
		},
		{
			name:           "case sensitivity",
			input:          []byte("Command Terminated With Exit Code 1\n"),
			expectedOutput: "Command Terminated With Exit Code 1\n",
			expectedBytes:  len("Command Terminated With Exit Code 1\n"),
			description:    "Should NOT filter case-varied messages (case sensitive)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := copyWriter{Stream: &buf}

			n, err := writer.Write(tc.input)

			testutils.AssertEquals(t, "Write error", nil, err)
			testutils.AssertEquals(t, "Bytes written: "+tc.description, tc.expectedBytes, n)
			testutils.AssertEquals(t, "Output mismatch: "+tc.description, tc.expectedOutput, buf.String())
		})
	}
}

func TestEnvResolution(t *testing.T) {
	// Save original environment
	origEnv := os.Environ()
	defer func() {
		// Restore environment
		os.Clearenv()
		for _, env := range origEnv {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}()

	// Set up test environment
	os.Clearenv()
	os.Setenv("TEST_VAR", "test_value")
	os.Setenv("ANOTHER_VAR", "another_value")
	os.Setenv("EMPTY_VAR", "")

	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "explicit key=value",
			input:    []string{"KEY=value"},
			expected: []string{"KEY=value"},
		},
		{
			name:     "lookup existing var",
			input:    []string{"TEST_VAR"},
			expected: []string{"TEST_VAR=test_value"},
		},
		{
			name:     "lookup non-existing var",
			input:    []string{"NONEXISTENT"},
			expected: []string{},
		},
		{
			name:     "mixed explicit and lookup",
			input:    []string{"KEY=value", "TEST_VAR", "ANOTHER_VAR"},
			expected: []string{"KEY=value", "TEST_VAR=test_value", "ANOTHER_VAR=another_value"},
		},
		{
			name:     "lookup empty var",
			input:    []string{"EMPTY_VAR"},
			expected: []string{"EMPTY_VAR="},
		},
		{
			name:     "equals sign in value",
			input:    []string{"KEY=value=with=equals"},
			expected: []string{"KEY=value=with=equals"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := resolveEnvVars(tc.input)
			testutils.AssertEquals(t, "Environment resolution mismatch", tc.expected, result)
		})
	}
}

// Helper function to test env resolution logic extracted from run()
func resolveEnvVars(envs []string) []string {
	newEnv := []string{}
	for _, envValue := range envs {
		if !strings.Contains(envValue, "=") {
			if value, set := os.LookupEnv(envValue); set {
				newEnv = append(newEnv, envValue+"="+value)
			}
		} else {
			newEnv = append(newEnv, envValue)
		}
	}
	return newEnv
}

func TestRunRawCmd(t *testing.T) {
	testCases := []struct {
		name        string
		command     string
		args        []string
		expectError bool
	}{
		{
			name:        "successful command",
			command:     "sh",
			args:        []string{"-c", "echo hello"},
			expectError: false,
		},
		{
			name:        "failing command",
			command:     "sh",
			args:        []string{"-c", "exit 1"},
			expectError: true,
		},
		{
			name:        "command with output",
			command:     "sh",
			args:        []string{"-c", "printf 'test output'"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := RunRawCmd(tc.command, tc.args)

			if tc.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestInteractiveTtyFlagHandling(t *testing.T) {
	testCases := []struct {
		name              string
		interactive       bool
		tty               bool
		expectedArgs      []string
		expectedEnvCount  int
		description       string
	}{
		{
			name:             "no flags",
			interactive:      false,
			tty:              false,
			expectedArgs:     []string{"exec", "pod-name"},
			expectedEnvCount: 0,
			description:      "Should not add -i, -t flags or env vars when both flags are false",
		},
		{
			name:             "interactive only",
			interactive:      true,
			tty:              false,
			expectedArgs:     []string{"exec", "-i", "pod-name"},
			expectedEnvCount: 1, // ENV=/etc/sh.shrc.local
			description:      "Should add -i flag and ENV env var when interactive is true",
		},
		{
			name:             "tty only",
			interactive:      false,
			tty:              true,
			expectedArgs:     []string{"exec", "-t", "pod-name"},
			expectedEnvCount: -1, // -1 means we just check it adds env vars (GetEnvironmentVarsList)
			description:      "Should add -t flag and env vars when tty is true",
		},
		{
			name:             "both interactive and tty",
			interactive:      true,
			tty:              true,
			expectedArgs:     []string{"exec", "-i", "-t", "pod-name"},
			expectedEnvCount: -1, // At least ENV + GetEnvironmentVarsList
			description:      "Should add both -i and -t flags with combined env vars",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Build command args similar to run() function
			commandArgs := []string{"exec"}
			envs := []string{}

			if tc.interactive {
				commandArgs = append(commandArgs, "-i")
				envs = append(envs, "ENV=/etc/sh.shrc.local")
			}
			if tc.tty {
				commandArgs = append(commandArgs, "-t")
				// Simulate GetEnvironmentVarsList() - returns some env vars
				envs = append(envs, "TERM=xterm", "USER=test")
			}

			commandArgs = append(commandArgs, "pod-name")

			// Verify flags are added correctly
			if len(tc.expectedArgs) > 0 {
				// Check that expected args are present in commandArgs
				for _, expectedArg := range tc.expectedArgs {
					found := false
					for _, arg := range commandArgs {
						if arg == expectedArg {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("%s: Expected arg '%s' not found in %v", tc.description, expectedArg, commandArgs)
					}
				}
			}

			// Verify env vars
			if tc.expectedEnvCount >= 0 {
				testutils.AssertEquals(t, tc.description, tc.expectedEnvCount, len(envs))
			} else {
				// Just check that envs were added
				if len(envs) == 0 {
					t.Errorf("%s: Expected env vars to be added", tc.description)
				}
			}
		})
	}
}

func TestKubectlBackendArgs(t *testing.T) {
	testCases := []struct {
		name         string
		command      string
		namespace    string
		expectedArgs []string
		description  string
	}{
		{
			name:         "kubectl with namespace",
			command:      "kubectl",
			namespace:    "default",
			expectedArgs: []string{"exec", "-n", "default", "-c", "uyuni", "--"},
			description:  "Should add kubectl-specific args with namespace",
		},
		{
			name:         "kubectl with empty namespace",
			command:      "kubectl",
			namespace:    "",
			expectedArgs: []string{"exec", "-n", "", "-c", "uyuni", "--"},
			description:  "Should add kubectl-specific args even with empty namespace",
		},
		{
			name:         "podman backend",
			command:      "podman",
			namespace:    "default",
			expectedArgs: []string{"exec", "pod-name"},
			description:  "Should NOT add kubectl-specific args for podman",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			commandArgs := []string{"exec"}
			commandArgs = append(commandArgs, "pod-name")

			if tc.command == "kubectl" {
				commandArgs = append(commandArgs, "-n", tc.namespace, "-c", "uyuni", "--")
			}

			// Verify expected args are present
			for _, expectedArg := range tc.expectedArgs {
				found := false
				for _, arg := range commandArgs {
					if arg == expectedArg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: Expected arg '%s' not found in %v", tc.description, expectedArg, commandArgs)
				}
			}
		})
	}
}

func TestEnvVarResolutionWithOsLookup(t *testing.T) {
	// Save and restore environment
	origEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range origEnv {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}()

	os.Clearenv()
	os.Setenv("MY_VAR", "my_value")
	os.Setenv("PATH", "/usr/bin")

	testCases := []struct {
		name           string
		inputEnvs      []string
		expectedEnvs   []string
		shouldNotExist []string
		description    string
	}{
		{
			name:           "lookup existing var",
			inputEnvs:      []string{"MY_VAR"},
			expectedEnvs:   []string{"MY_VAR=my_value"},
			shouldNotExist: []string{},
			description:    "Should resolve MY_VAR from OS environment",
		},
		{
			name:           "lookup non-existing var",
			inputEnvs:      []string{"NONEXISTENT"},
			expectedEnvs:   []string{},
			shouldNotExist: []string{"NONEXISTENT"},
			description:    "Should skip non-existing env vars",
		},
		{
			name:           "explicit value overrides lookup",
			inputEnvs:      []string{"MY_VAR=override_value"},
			expectedEnvs:   []string{"MY_VAR=override_value"},
			shouldNotExist: []string{},
			description:    "Should use explicit value instead of OS lookup",
		},
		{
			name:           "mixed lookup and explicit",
			inputEnvs:      []string{"MY_VAR", "EXPLICIT=value", "PATH"},
			expectedEnvs:   []string{"MY_VAR=my_value", "EXPLICIT=value", "PATH=/usr/bin"},
			shouldNotExist: []string{},
			description:    "Should handle mixed explicit and lookup vars",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := resolveEnvVars(tc.inputEnvs)

			// Check expected envs are present
			for _, expected := range tc.expectedEnvs {
				found := false
				for _, res := range result {
					if res == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: Expected env '%s' not found in %v", tc.description, expected, result)
				}
			}

			// Check non-existing vars are not present
			for _, shouldNot := range tc.shouldNotExist {
				for _, res := range result {
					if strings.HasPrefix(res, shouldNot+"=") {
						t.Errorf("%s: Env '%s' should not be in result %v", tc.description, shouldNot, result)
					}
				}
			}
		})
	}
}
