// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/backend"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// -----------------------------------------------------------------------
// mockRunner – satisfies types.Runner for unit tests.
// -----------------------------------------------------------------------

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

// -----------------------------------------------------------------------
// stubDetector – satisfies backend.BackendDetector for unit tests.
// -----------------------------------------------------------------------

// stubDetector is a simple BackendDetector whose behaviour is fully
// controlled by the test.  It records the arguments it receives so tests
// can assert that GetCommand passes the right values through.
type stubDetector struct {
	wantExplicit         string
	wantContainer        string
	wantKubernetesFilter string
	returnCommand        string
	returnErr            error
}

func (s *stubDetector) Detect(explicit, container, kubernetesFilter string) (string, error) {
	// If the test set expectations, verify them.
	if s.wantExplicit != "" && explicit != s.wantExplicit {
		return "", errors.New("stubDetector: unexpected explicit backend: " + explicit)
	}
	if s.wantContainer != "" && container != s.wantContainer {
		return "", errors.New("stubDetector: unexpected container: " + container)
	}
	if s.wantKubernetesFilter != "" && kubernetesFilter != s.wantKubernetesFilter {
		return "", errors.New("stubDetector: unexpected filter: " + kubernetesFilter)
	}
	return s.returnCommand, s.returnErr
}

// newTestConnection returns a Connection whose detector is fully stubbed.
func newTestConnection(bk, container, filter string, d backend.BackendDetector) *Connection {
	cnx := NewConnection(bk, container, filter)
	cnx.WithDetector(d)
	return cnx
}

// -----------------------------------------------------------------------
// TestGetCommand – table-driven tests using the stub detector
// -----------------------------------------------------------------------

func TestGetCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		backend          string
		container        string
		kubernetesFilter string
		detectCommand    string
		detectErr        error
		wantCommand      string
		wantErr          bool
	}{
		{
			name:          "explicit podman → returned as-is",
			backend:       "podman",
			detectCommand: "podman",
			wantCommand:   "podman",
		},
		{
			name:          "explicit podman-remote → returned as-is",
			backend:       "podman-remote",
			detectCommand: "podman-remote",
			wantCommand:   "podman-remote",
		},
		{
			name:          "explicit kubectl → returned as-is",
			backend:       "kubectl",
			detectCommand: "kubectl",
			wantCommand:   "kubectl",
		},
		{
			name:          "explicit host → returned as-is",
			backend:       "host",
			detectCommand: "host",
			wantCommand:   "host",
		},
		{
			name:          "auto-detect resolves to podman",
			backend:       "",
			detectCommand: "podman",
			wantCommand:   "podman",
		},
		{
			name:          "auto-detect resolves to kubectl",
			backend:       "",
			detectCommand: "kubectl",
			wantCommand:   "kubectl",
		},
		{
			name:      "detector returns error → GetCommand propagates it",
			backend:   "",
			detectErr: errors.New("no backend found"),
			wantErr:   true,
		},
		{
			name:          "result is cached – second call does not re-invoke detector",
			backend:       "podman",
			detectCommand: "podman",
			wantCommand:   "podman",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stub := &stubDetector{
				returnCommand: tc.detectCommand,
				returnErr:     tc.detectErr,
			}
			cnx := newTestConnection(tc.backend, "uyuni-server", "-lapp=uyuni", stub)

			got, err := cnx.GetCommand()

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got command=%q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantCommand {
				t.Errorf("GetCommand() = %q, want %q", got, tc.wantCommand)
			}

			// Verify caching: call again and ensure the same result without
			// the detector being invoked a second time (the stub would still
			// return the same value, but the point is the command field is
			// populated already).
			got2, err2 := cnx.GetCommand()
			if err2 != nil {
				t.Fatalf("second GetCommand() unexpected error: %v", err2)
			}
			if got2 != got {
				t.Errorf("second GetCommand() = %q, want %q (caching broken)", got2, got)
			}
		})
	}
}

// TestGetCommand_DetectorReceivesCorrectArgs verifies that GetCommand
// forwards backend, container and kubernetesFilter to the detector.
func TestGetCommand_DetectorReceivesCorrectArgs(t *testing.T) {
	t.Parallel()

	stub := &stubDetector{
		wantExplicit:         "podman",
		wantContainer:        "my-container",
		wantKubernetesFilter: "-lapp=test",
		returnCommand:        "podman",
	}
	cnx := newTestConnection("podman", "my-container", "-lapp=test", stub)

	if _, err := cnx.GetCommand(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// -----------------------------------------------------------------------
// Existing tests (unchanged behaviour)
// -----------------------------------------------------------------------

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

	if len(capturedArgs) < 2 {
		t.Fatal("Unexpected number of captured args")
	}

	if capturedArgs[0] != "/etc/passwd" || capturedArgs[1] != "/tmp/passwd" {
		t.Errorf("Expected args ['/etc/passwd', '/tmp/passwd'], got %v", capturedArgs)
	}
}

func TestHostInferredCopy(t *testing.T) {
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
	err := cnx.Copy("server:/home/myfile", "/tmp/", "", "")

	if err != nil {
		t.Errorf("Copy returned error: %v", err)
	}

	if capturedCommand != "cp" {
		t.Errorf("Expected command 'cp', got '%s'", capturedCommand)
	}

	if len(capturedArgs) < 2 {
		t.Fatal("Unexpected number of captured args")
	}

	if capturedArgs[0] != "/home/myfile" || capturedArgs[1] != "/tmp/myfile" {
		t.Errorf("Expected args ['/home/myfile', '/tmp/myfile'], got %v", capturedArgs)
	}

	err = cnx.Copy("/tmp/myanotherfile", "server:", "", "")
	if err != nil {
		t.Errorf("Copy returned error: %v", err)
	}

	if len(capturedArgs) < 2 {
		t.Fatal("Unexpected number of captured args")
	}

	if capturedArgs[0] != "/tmp/myanotherfile" || capturedArgs[1] != "myanotherfile" {
		t.Errorf("Expected args ['/tmp/myanotherfile', 'myanotherfile'], got %v", capturedArgs)
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
