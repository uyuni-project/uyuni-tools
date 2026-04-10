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
