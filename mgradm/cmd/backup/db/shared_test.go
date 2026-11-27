// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParsePostgresConfig(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "postgresql.conf")
	content := `
# This is a comment
listen_addresses = '*' # inline comment
port = 5432
max_connections = 100
# archive_mode = off
archive_mode = on
archive_command = '/usr/bin/smdba-pgarchive --source "%p" --destination "/var/lib/pgsql/backup/%f"'
restore_command = '/usr/bin/cp /var/lib/pgsql/backup/%f %p'
`
	testutils.WriteFile(t, configPath, content)

	podman.SetRunner(func(command string, args ...string) types.Runner {
		if command == "podman" && args[0] == "volume" && args[1] == "inspect" {
			return testutils.FakeRunnerGenerator(tmpDir, nil)(command, args...)
		}
		return testutils.FakeRunnerGenerator("", nil)(command, args...)
	})
	defer podman.ResetRunner()

	config, err := ParsePostgresConfig()
	testutils.AssertNoError(t, "ParsePostgresConfig failed", err)

	expected := map[string]string{
		"listen_addresses": "'*'",
		"port":             "5432",
		"max_connections":  "100",
		"archive_mode":     "on",
		"archive_command":  "'/usr/bin/smdba-pgarchive --source \"%p\" --destination \"/var/lib/pgsql/backup/%f\"'",
		"restore_command":  "'/usr/bin/cp /var/lib/pgsql/backup/%f %p'",
	}

	testutils.AssertEquals(t, "config mismatch", expected, config)
}

func TestUpdatePostgresConfig(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "postgresql.conf")
	initialContent := `# Initial config
listen_addresses = '*'
port = 5432
`
	testutils.WriteFile(t, configPath, initialContent)

	podman.SetRunner(func(command string, args ...string) types.Runner {
		if command == "podman" && args[0] == "volume" && args[1] == "inspect" {
			return testutils.FakeRunnerGenerator(tmpDir, nil)(command, args...)
		}
		return testutils.FakeRunnerGenerator("", nil)(command, args...)
	})
	defer podman.ResetRunner()

	updates := map[string]string{
		"port":            "5433",
		"archive_mode":    "on",
		"archive_command": "'/usr/bin/smdba-pgarchive --source \"%p\" --destination \"/var/lib/pgsql/backup/%f\"'",
	}

	err := UpdatePostgresConfig(updates)
	testutils.AssertNoError(t, "UpdatePostgresConfig failed", err)

	// Verify the file content indirectly via ParsePostgresConfig
	config, err := ParsePostgresConfig()
	testutils.AssertNoError(t, "ParsePostgresConfig failed", err)

	expected := map[string]string{
		"listen_addresses": "'*'",
		"port":             "5433",
		"archive_mode":     "on",
		"archive_command":  "'/usr/bin/smdba-pgarchive --source \"%p\" --destination \"/var/lib/pgsql/backup/%f\"'",
	}

	testutils.AssertEquals(t, "config mismatch", expected, config)

	// Verify that comments are preserved and format is correct
	updatedContent := testutils.ReadFile(t, configPath)
	expectedFileContent := `# Initial config
listen_addresses = '*'
port = 5433
archive_mode = on
archive_command = '/usr/bin/smdba-pgarchive --source "%p" --destination "/var/lib/pgsql/backup/%f"'
`
	testutils.AssertTrue(t, "Unexpected updated content:\n"+testutils.DiffStrings(expectedFileContent, updatedContent),
		areStringsSame(updatedContent, expectedFileContent))
}

func areStringsSame(s1, s2 string) bool {
	lines1 := strings.Split(s1, "\n")
	lines2 := strings.Split(s2, "\n")

	// If the total number of lines differs, they can't be identical
	if len(lines1) != len(lines2) {
		return false
	}

	counts := make(map[string]int)

	for _, line := range lines1 {
		counts[line]++
	}

	for _, line := range lines2 {
		if counts[line] == 0 {
			// Line doesn't exist in s1 or appears more often in s2
			return false
		}
		counts[line]--
	}

	// Everything in those lines should match, count should be 0
	for _, count := range counts {
		if count != 0 {
			return false
		}
	}
	return true
}
