// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

type mockRunner struct {
	output []byte
	err    error
}

func (m *mockRunner) Log(_ zerolog.Level) types.Runner  { return m }
func (m *mockRunner) Spinner(_ string) types.Runner     { return m }
func (m *mockRunner) StdMapping() types.Runner          { return m }
func (m *mockRunner) Std(_ *bytes.Buffer) types.Runner  { return m }
func (m *mockRunner) Wait() error                       { return nil }
func (m *mockRunner) InputString(_ string) types.Runner { return m }
func (m *mockRunner) Env(_ []string) types.Runner       { return m }
func (m *mockRunner) Start() error                      { return nil }

func (m *mockRunner) Exec() ([]byte, error) {
	return m.output, m.err
}

func TestExecScript(t *testing.T) {
	originalRunner := runner
	defer func() { runner = originalRunner }()

	var capturedCommands [][]string

	runner = func(command string, args ...string) types.Runner {
		capturedCommands = append(capturedCommands, append([]string{command}, args...))
		if command == "podman" {
			// exec command
			if len(args) > 0 && args[0] == "exec" {
				return &mockRunner{output: []byte("script output")}
			}
		}
		return &mockRunner{}
	}

	cnx := NewConnection("podman", "uyuni-server", "")
	cnx.command = "podman"
	cnx.podName = "uyuni-server"

	scriptContent := "echo 'hello'"
	out, err := cnx.ExecScript(scriptContent)

	if err != nil {
		t.Errorf("ExecScript returned error: %v", err)
	}

	if string(out) != "script output" {
		t.Errorf("Expected output 'script output', got '%s'", string(out))
	}

	foundCp := false
	foundExec := false
	foundRm := false

	for _, cmd := range capturedCommands {
		if len(cmd) <= 1 {
			// ignore wrong commands, will fail later if it is a problem
			continue
		}
		if cmd[0] == "podman" {
			switch cmd[1] {
			case "cp":
				foundCp = true
				if !strings.Contains(cmd[3], "server:/tmp/script-") {
					t.Errorf("Unexpected cp destination: %s", cmd[3])
				}
			case "exec":
				// cmd[2] is a container name
				// Check if it is the bash execution
				if cmd[3] == "bash" && strings.Contains(strings.Join(cmd[4:], " "), "/tmp/script-") {
					foundExec = true
				}
				// Check if it is the cleanup
				if cmd[3] == "rm" && strings.Contains(strings.Join(cmd[4:], " "), "-f /tmp/script-") {
					foundRm = true
				}
			}
		}
	}

	if !foundCp {
		t.Error("Did not find copy command")
	}
	if !foundExec {
		t.Error("Did not find execution command")
	}
	if !foundRm {
		t.Error("Did not find cleanup command")
	}
}

func TestHostExec(t *testing.T) {
	originalRunner := runner
	defer func() { runner = originalRunner }()

	var capturedCommand string
	var capturedArgs []string

	runner = func(command string, args ...string) types.Runner {
		capturedCommand = command
		capturedArgs = args
		return &mockRunner{output: []byte("host output")}
	}

	cnx := NewConnection("host", "", "")
	out, err := cnx.Exec("ls", "-la")

	if err != nil {
		t.Errorf("Exec returned error: %v", err)
	}

	if capturedCommand != "ls" {
		t.Errorf("Expected command 'ls', got '%s'", capturedCommand)
	}

	if len(capturedArgs) != 1 || capturedArgs[0] != "-la" {
		t.Errorf("Expected args ['-la'], got %v", capturedArgs)
	}

	if string(out) != "host output" {
		t.Errorf("Expected output 'host output', got '%s'", string(out))
	}
}

func TestHostCopy(t *testing.T) {
	originalRunner := runner
	defer func() { runner = originalRunner }()

	var capturedCommand string
	var capturedArgs []string

	runner = func(command string, args ...string) types.Runner {
		capturedCommand = command
		capturedArgs = args
		return &mockRunner{}
	}

	cnx := NewConnection("host", "", "")
	err := cnx.Copy("server:/etc/passwd", "/tmp/passwd", "", "")

	if err != nil {
		t.Errorf("Copy returned error: %v", err)
	}

	if capturedCommand != "cp" {
		t.Errorf("Expected command 'cp', got '%s'", capturedCommand)
	}

	if capturedArgs[0] != "/etc/passwd" || capturedArgs[1] != "/tmp/passwd" {
		t.Errorf("Expected args ['/etc/passwd', '/tmp/passwd'], got %v", capturedArgs)
	}
}

func TestHostUserExec(t *testing.T) {
	originalRunner := runner
	defer func() { runner = originalRunner }()

	var capturedCommand string
	var capturedArgs []string

	runner = func(command string, args ...string) types.Runner {
		capturedCommand = command
		capturedArgs = args
		return &mockRunner{output: []byte("host output")}
	}

	cnx := NewUserConnection("host", "", "", "postgres")
	out, err := cnx.Exec("ls", "-la")

	if err != nil {
		t.Errorf("Exec returned error: %v", err)
	}

	if capturedCommand != "su" {
		t.Errorf("Expected command 'su', got '%s'", capturedCommand)
	}

	expectedArgs := []string{"-", "postgres", "-c", "'ls' '-la'"}
	if len(capturedArgs) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d: %v", len(expectedArgs), len(capturedArgs), capturedArgs)
	} else {
		for i, arg := range capturedArgs {
			if arg != expectedArgs[i] {
				t.Errorf("Expected arg %d to be '%s', got '%s'", i, expectedArgs[i], arg)
			}
		}
	}

	if string(out) != "host output" {
		t.Errorf("Expected output 'host output', got '%s'", string(out))
	}
}

func TestPodmanUserExec(t *testing.T) {
	originalRunner := runner
	defer func() { runner = originalRunner }()

	var capturedCommand string
	var capturedArgs []string

	runner = func(command string, args ...string) types.Runner {
		capturedCommand = command
		capturedArgs = args
		return &mockRunner{output: []byte("podman output")}
	}

	cnx := NewUserConnection("podman", "uyuni-server", "", "postgres")
	cnx.command = "podman"
	cnx.podName = "uyuni-server"
	out, err := cnx.Exec("psql", "-c", "SHOW archive_mode;")

	if err != nil {
		t.Errorf("Exec returned error: %v", err)
	}

	if capturedCommand != "podman" {
		t.Errorf("Expected command 'podman', got '%s'", capturedCommand)
	}

	expectedArgs := []string{"exec", "uyuni-server", "su", "-", "postgres", "-c", "'psql' '-c' 'SHOW archive_mode;'"}
	if len(capturedArgs) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d: %v", len(expectedArgs), len(capturedArgs), capturedArgs)
	} else {
		for i, arg := range capturedArgs {
			if arg != expectedArgs[i] {
				t.Errorf("Expected arg %d to be '%s', got '%s'", i, expectedArgs[i], arg)
			}
		}
	}

	if string(out) != "podman output" {
		t.Errorf("Expected output 'podman output', got '%s'", string(out))
	}
}

func TestQuoteArgs(t *testing.T) {
	tests := []struct {
		args     []string
		expected string
	}{
		{args: []string{"ls", "-la"}, expected: "'ls' '-la'"},
		{args: []string{"echo", "hello world"}, expected: "'echo' 'hello world'"},
		{args: []string{"psql", "-c", "SHOW archive_mode;"}, expected: "'psql' '-c' 'SHOW archive_mode;'"},
		{args: []string{"sh", "-c", "echo 'hello'"}, expected: "'sh' '-c' 'echo '\\''hello'\\'''"},
	}

	for _, tt := range tests {
		actual := quoteArgs(tt.args)
		if actual != tt.expected {
			t.Errorf("quoteArgs(%v) = %s, expected %s", tt.args, actual, tt.expected)
		}
	}
}
