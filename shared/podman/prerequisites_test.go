// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os/exec"
	"testing"
)

func TestCheckPodmanRunningContainers(t *testing.T) {
	// Skip test if podman is not installed
	if _, err := exec.LookPath("podman"); err != nil {
		t.Skip("podman not installed, skipping test")
	}

	err := CheckPodmanRunningContainers()
	// Should not error by default if there are no uyuni containers
	if err != nil {
		t.Logf("Expected no error or error if uyuni network has containers, got: %s", err)
	}
}
