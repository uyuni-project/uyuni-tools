// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestCheckPodmanRunningContainers(t *testing.T) {
	type testCase struct {
		name      string
		output    string
		err       error
		wantError string
	}

	tests := []testCase{
		{
			name: "no containers running",
		},
		{
			name:      "containers running",
			output:    "abc123\n",
			wantError: "running containers",
		},
		{
			name:      "podman command fails",
			err:       errors.New("podman unavailable"),
			wantError: "failed to check running podman containers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldRunner := runner
			defer func() { runner = oldRunner }()

			var command string
			var args []string
			runner = func(gotCommand string, gotArgs ...string) types.Runner {
				command = gotCommand
				args = gotArgs
				return testutils.FakeRunnerGenerator(tt.output, tt.err)(gotCommand, gotArgs...)
			}

			err := CheckPodmanRunningContainers("test-network")
			if tt.wantError == "" && err != nil {
				t.Fatalf("CheckPodmanRunningContainers() unexpected error: %v", err)
			}
			if tt.wantError != "" && (err == nil || !strings.Contains(err.Error(), tt.wantError)) {
				t.Fatalf("CheckPodmanRunningContainers() error = %v, want message containing %q", err, tt.wantError)
			}
			if command != "podman" || strings.Join(args, " ") != "ps -q --filter network=test-network" {
				t.Errorf("runner called with %q %q", command, args)
			}
		})
	}
}

func TestIsCheckSkipped(t *testing.T) {
	for _, value := range []string{"1", "true", "TRUE"} {
		t.Setenv(skipPrerequisitesEnv, value)
		if !IsCheckSkipped() {
			t.Errorf("IsCheckSkipped() = false for %q", value)
		}
	}

	for _, value := range []string{"", "0", "false", "yes"} {
		t.Setenv(skipPrerequisitesEnv, value)
		if IsCheckSkipped() {
			t.Errorf("IsCheckSkipped() = true for %q", value)
		}
	}
}
