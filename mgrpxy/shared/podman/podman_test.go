// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"path"
	"testing"
)

func TestCheckDirPermissions(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := checkPermissions(tempDir, 0o005|0o050|0o500); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateYamlFiles(t *testing.T) {
	tempDir := t.TempDir()
	testFiles := []string{"httpd.yaml", "ssh.yaml", "config.yaml"}
	for _, file := range testFiles {
		filePath := path.Join(tempDir, file)
		if _, err := os.Create(filePath); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filePath, err)
		}
	}

	// Test: when all files are present and have correct permissions
	if err := validateInstallYamlFiles(tempDir); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Change the permission of config.yaml to 0600 to simulate a permission error
	configFilePath := path.Join(tempDir, "config.yaml")
	if err := os.Chmod(configFilePath, 0o600); err != nil {
		t.Fatalf("Failed to change permissions for %s: %v", configFilePath, err)
	}
	if err := validateInstallYamlFiles(tempDir); err == nil {
		t.Errorf("Expected an error due to incorrect permissions on config.yaml, but got none")
	}

	// Restore the correct permissions for the next test run
	if err := os.Chmod(configFilePath, 0o644); err != nil {
		t.Fatalf("Failed to restore permissions for %s: %v", configFilePath, err)
	}

	// Test: Missing file scenario, remove one file and expect an error
	os.Remove(path.Join(tempDir, "httpd.yaml"))
	if err := validateInstallYamlFiles(tempDir); err == nil {
		t.Errorf("Expected an error due to missing httpd.yaml, but got none")
	}
}
