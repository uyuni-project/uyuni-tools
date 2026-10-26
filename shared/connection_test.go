// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
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

const testHealthcheckConfig = `{"StartPeriod":600000000000,"Interval":10000000000,"Timeout":120000000000,"Retries":1}`

func setRunnerHandler(t *testing.T, handler func(command string, args ...string) ([]byte, error)) {
	t.Helper()
	originalRunner := runner
	t.Cleanup(func() { runner = originalRunner })

	runner = func(command string, args ...string) types.Runner {
		output, err := handler(command, args...)
		return &mockRunner{output: output, err: err}
	}
	podman.SetRunner(runner)
	t.Cleanup(podman.ResetRunner)
}

func setShortWaitTimeouts(t *testing.T) {
	t.Helper()
	appearanceTimeout := containerAppearanceTimeout
	noEventTimeout := containerNoEventTimeout
	pollInterval := healthcheckPollInterval
	containerAppearanceTimeout = 5 * time.Millisecond
	containerNoEventTimeout = 2 * time.Millisecond
	healthcheckPollInterval = time.Millisecond
	t.Cleanup(func() {
		containerAppearanceTimeout = appearanceTimeout
		containerNoEventTimeout = noEventTimeout
		healthcheckPollInterval = pollInterval
	})
}

func testEvent(id string, status string, timeNano int64, exitCode int) string {
	return fmt.Sprintf(`{"ID":%q,"Status":%q,"timeNano":%d,"ContainerExitCode":%d}`,
		id, status, timeNano, exitCode)
}

func testEvents(events ...string) []byte {
	return []byte(strings.Join(events, "\n"))
}

func testInspectOutput(
	state string, running bool, exitCode int, health string, failingStreak int, healthLog string,
) []byte {
	healthJSON := "null"
	if health != "" {
		healthJSON = fmt.Sprintf(`{"Status":%q,"FailingStreak":%d,"Log":[{"ExitCode":0,"Output":%q}]}`,
			health, failingStreak, healthLog)
	}

	return []byte(fmt.Sprintf(
		`{"RestartCount":0,"Config":{"Healthcheck":%s},`+
			`"State":{"Status":%q,"Running":%t,"ExitCode":%d,`+
			`"StartedAt":"2026-10-23T12:00:00Z","Health":%s}}`,
		testHealthcheckConfig, state, running, exitCode, healthJSON,
	))
}

func testRunner(
	events []byte, eventsErr error, inspect []byte, inspectErr error,
) func(command string, args ...string) ([]byte, error) {
	return func(command string, args ...string) ([]byte, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("unexpected command: %s", command)
		}
		switch args[0] {
		case "events":
			return events, eventsErr
		case "inspect":
			return inspect, inspectErr
		default:
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		}
	}
}

func TestWaitForHealthcheck(t *testing.T) {
	now := time.Now()
	healthyEvents := testEvents(
		testEvent("container-id", "create", now.Add(time.Millisecond).UnixNano(), 0),
		testEvent("container-id", "start", now.Add(2*time.Millisecond).UnixNano(), 0),
	)
	healthyInspect := testInspectOutput("running", true, 0, "healthy", 0, "")

	tests := []struct {
		name       string
		events     []byte
		inspect    []byte
		inspectErr error
		wantErr    string
	}{
		{
			name:    "healthy",
			events:  healthyEvents,
			inspect: healthyInspect,
		},
		{
			name: "multiple starts",
			events: testEvents(
				testEvent("container-id-1", "start", now.Add(time.Millisecond).UnixNano(), 0),
				testEvent("container-id-2", "start", now.Add(2*time.Millisecond).UnixNano(), 0),
			),
			inspect: healthyInspect,
			wantErr: "more than once",
		},
		{
			name: "exited",
			events: testEvents(
				testEvent("container-id", "create", now.Add(time.Millisecond).UnixNano(), 0),
				testEvent("container-id", "start", now.Add(2*time.Millisecond).UnixNano(), 0),
				testEvent("container-id", "died", now.Add(3*time.Millisecond).UnixNano(), 42),
			),
			inspectErr: errors.New("no such object"),
			wantErr:    "exit code 42",
		},
		{
			name: "removed before start",
			events: testEvents(
				testEvent("container-id", "create", now.Add(time.Millisecond).UnixNano(), 0),
				testEvent("container-id", "remove", now.Add(2*time.Millisecond).UnixNano(), 0),
			),
			inspectErr: errors.New("no such object"),
			wantErr:    "before it started",
		},
		{
			name:    "no healthcheck",
			events:  healthyEvents,
			inspect: testInspectOutput("running", true, 0, "", 0, ""),
			wantErr: "no healthcheck",
		},
		{
			name:    "unhealthy",
			events:  healthyEvents,
			inspect: testInspectOutput("running", true, 0, "unhealthy", 3, "healthcheck failed"),
			wantErr: "healthcheck failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setRunnerHandler(t, testRunner(tt.events, nil, tt.inspect, tt.inspectErr))

			cnx := NewConnection("podman", "uyuni-server", "")
			cnx.command = "podman"

			err := cnx.WaitForHealthcheck(now)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected an error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestWaitForHealthcheckKubernetesUnsupported(t *testing.T) {
	cnx := NewConnection("kubectl", "uyuni-server", "")
	cnx.command = "kubectl"

	err := cnx.WaitForHealthcheck(time.Now())
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("expected a not supported error, got: %v", err)
	}
}

func TestWaitForHealthcheckHost(t *testing.T) {
	cnx := NewConnection("host", "", "")
	cnx.command = "host"

	if err := cnx.WaitForHealthcheck(time.Now()); err != nil {
		t.Fatalf("unexpected error for the host backend: %v", err)
	}
}

func TestWaitForHealthcheckNoEventsTimeout(t *testing.T) {
	setShortWaitTimeouts(t)
	setRunnerHandler(t, testRunner(nil, nil, nil, errors.New("no such object")))

	cnx := NewConnection("podman", "uyuni-server", "")
	cnx.command = "podman"

	start := time.Now()
	err := cnx.WaitForHealthcheck(start)
	if err == nil || !strings.Contains(err.Error(), "probably failed to start it") {
		t.Fatalf("expected a failed to start error, got: %v", err)
	}
	// The no event timeout must be used instead of the longer appearance timeout.
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("expected the no event timeout to be used, waited %s", elapsed)
	}
}

func TestWaitForHealthcheckAppearanceTimeout(t *testing.T) {
	setShortWaitTimeouts(t)
	now := time.Now()
	// The container got created but never became inspectable.
	events := testEvents(testEvent("container-id", "create", now.Add(time.Millisecond).UnixNano(), 0))
	setRunnerHandler(t, testRunner(events, nil, nil, errors.New("no such object")))

	cnx := NewConnection("podman", "uyuni-server", "")
	cnx.command = "podman"

	err := cnx.WaitForHealthcheck(now)
	if err == nil || !strings.Contains(err.Error(), "did not appear") {
		t.Fatalf("expected a did not appear error, got: %v", err)
	}
}

func TestWaitForHealthcheckCreatedButNotStarted(t *testing.T) {
	setShortWaitTimeouts(t)
	now := time.Now()
	createEvent := testEvents(testEvent("container-id", "create", now.Add(time.Millisecond).UnixNano(), 0))
	inspect := testInspectOutput("created", false, 0, "starting", 0, "")
	setRunnerHandler(t, testRunner(createEvent, nil, inspect, nil))

	cnx := NewConnection("podman", "uyuni-server", "")
	cnx.command = "podman"

	err := cnx.WaitForHealthcheck(now)
	if err == nil || !strings.Contains(err.Error(), "did not start") {
		t.Fatalf("expected a did not start error, got: %v", err)
	}
}

func TestWaitForHealthcheckEventsUnavailable(t *testing.T) {
	setRunnerHandler(t, testRunner(nil, errors.New("events unavailable"), nil, nil))

	cnx := NewConnection("podman", "uyuni-server", "")
	cnx.command = "podman"

	err := cnx.WaitForHealthcheck(time.Now())
	if err == nil || !strings.Contains(err.Error(), "cannot read the startup events") {
		t.Fatalf("expected an events error, got: %v", err)
	}
}

func TestGetPodNameUsesExactFilter(t *testing.T) {
	var filter string
	setRunnerHandler(t, func(_ string, args ...string) ([]byte, error) {
		if len(args) >= 4 && args[0] == "ps" {
			filter = args[3]
			return []byte("container-id\n"), nil
		}
		return nil, fmt.Errorf("unexpected command: %v", args)
	})

	cnx := NewConnection("podman", "uyuni-server", "")
	cnx.command = "podman"

	if _, err := cnx.GetPodName(); err != nil {
		t.Fatalf("GetPodName returned an error: %v", err)
	}
	if filter != "name=^uyuni-server$" {
		t.Fatalf("unexpected name filter: %s", filter)
	}
}

func TestWaitForContainer(t *testing.T) {
	now := time.Now()
	runningInspect := testInspectOutput("running", true, 0, "healthy", 0, "")
	healthyEvents := testEvents(
		testEvent("container-id", "create", now.Add(time.Millisecond).UnixNano(), 0),
		testEvent("container-id", "start", now.Add(2*time.Millisecond).UnixNano(), 0),
	)

	tests := []struct {
		name       string
		events     []byte
		inspect    []byte
		inspectErr error
		wantErr    string
	}{
		{
			name:    "running",
			events:  healthyEvents,
			inspect: runningInspect,
		},
		{
			name: "old container death ignored on restart",
			events: testEvents(
				testEvent("old-id", "died", now.Add(time.Millisecond).UnixNano(), 0),
				testEvent("old-id", "remove", now.Add(2*time.Millisecond).UnixNano(), 0),
				testEvent("new-id", "create", now.Add(3*time.Millisecond).UnixNano(), 0),
				testEvent("new-id", "start", now.Add(4*time.Millisecond).UnixNano(), 0),
			),
			inspect: runningInspect,
		},
		{
			name: "multiple starts",
			events: testEvents(
				testEvent("container-id-1", "create", now.Add(time.Millisecond).UnixNano(), 0),
				testEvent("container-id-1", "start", now.Add(2*time.Millisecond).UnixNano(), 0),
				testEvent("container-id-2", "create", now.Add(3*time.Millisecond).UnixNano(), 0),
				testEvent("container-id-2", "start", now.Add(4*time.Millisecond).UnixNano(), 0),
			),
			inspect: runningInspect,
			wantErr: "more than once",
		},
		{
			name: "new container died",
			events: testEvents(
				testEvent("container-id", "create", now.Add(time.Millisecond).UnixNano(), 0),
				testEvent("container-id", "start", now.Add(2*time.Millisecond).UnixNano(), 0),
				testEvent("container-id", "died", now.Add(3*time.Millisecond).UnixNano(), 42),
			),
			inspectErr: errors.New("no such object"),
			wantErr:    "exit code 42",
		},
		{
			name:    "exited",
			inspect: testInspectOutput("exited", false, 1, "starting", 0, ""),
			wantErr: "not running",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setRunnerHandler(t, testRunner(tt.events, nil, tt.inspect, tt.inspectErr))

			cnx := NewConnection("podman", "uyuni-server", "")
			cnx.command = "podman"

			err := cnx.WaitForContainer(now)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected an error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestWaitForContainerNoEventsTimeout(t *testing.T) {
	setShortWaitTimeouts(t)
	setRunnerHandler(t, testRunner(nil, nil, nil, errors.New("no such object")))

	cnx := NewConnection("podman", "uyuni-server", "")
	cnx.command = "podman"

	start := time.Now()
	err := cnx.WaitForContainer(start)
	if err == nil || !strings.Contains(err.Error(), "probably failed to start it") {
		t.Fatalf("expected a failed to start error, got: %v", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("expected the no event timeout to be used, waited %s", elapsed)
	}
}

func TestWaitForContainerAppearanceTimeout(t *testing.T) {
	setShortWaitTimeouts(t)
	now := time.Now()
	// The container got created but never became inspectable.
	events := testEvents(testEvent("container-id", "create", now.Add(time.Millisecond).UnixNano(), 0))
	setRunnerHandler(t, testRunner(events, nil, nil, errors.New("no such object")))

	cnx := NewConnection("podman", "uyuni-server", "")
	cnx.command = "podman"

	err := cnx.WaitForContainer(now)
	if err == nil || !strings.Contains(err.Error(), "did not appear") {
		t.Fatalf("expected a did not appear error, got: %v", err)
	}
}
