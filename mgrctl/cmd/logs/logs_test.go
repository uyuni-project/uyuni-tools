// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package logs

import (
	"strings"
	"testing"
)

func TestValidateFlags(t *testing.T) {
	// Test --files without --reposync (Should Fail)
	invalidFlags := &flagpole{Files: "myregex"}
	err := validateFlags(invalidFlags)
	if err == nil {
		t.Errorf("Expected error when using --files without --reposync, got nil")
	}

	// Test --files with --reposync (Should Pass)
	validFlags := &flagpole{Reposync: true, Files: "myregex"}
	err = validateFlags(validFlags)
	if err != nil {
		t.Errorf("Expected no error when using --files with --reposync, got %v", err)
	}
}

func TestGetLogPaths(t *testing.T) {
	flags := &flagpole{
		Salt: true,
	}

	paths := getLogPaths(flags)

	// Test 1: Ensure salt paths exist but minion is absent
	pathStr := strings.Join(paths, " ")
	if !strings.Contains(pathStr, "/var/log/salt/api") || !strings.Contains(pathStr, "/var/log/salt/master") {
		t.Errorf("Expected salt api and master paths, got: %v", paths)
	}
	if strings.Contains(pathStr, "/var/log/salt/minion") {
		t.Errorf("Did not expect minion path in salt logs, got: %v", paths)
	}

	// Test 2: Ensure empty paths if no flags are set
	emptyFlags := &flagpole{}
	emptyPaths := getLogPaths(emptyFlags)
	if len(emptyPaths) != 0 {
		t.Errorf("Expected 0 paths for empty flags, got: %v", emptyPaths)
	}
}
