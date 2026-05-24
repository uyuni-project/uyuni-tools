// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestCheckPodmanRunningContainersNoContainers(t *testing.T) {
	// Mock runner to simulate no running containers
	oldRunner := runner
	defer func() { runner = oldRunner }()

	runner = testutils.FakeRunnerGenerator("", nil)

	err := CheckPodmanRunningContainers()
	if err != nil {
		t.Errorf("Expected no error when no containers are running, got: %v", err)
	}
}

func TestCheckPodmanRunningContainersWithContainers(t *testing.T) {
	// Mock runner to simulate running containers
	oldRunner := runner
	defer func() { runner = oldRunner }()

	runner = testutils.FakeRunnerGenerator("abc123\n", nil)

	err := CheckPodmanRunningContainers()
	if err == nil {
		t.Error("Expected error when containers are running, got nil")
	}
}
