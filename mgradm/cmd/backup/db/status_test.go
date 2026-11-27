// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/pgsql"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestStatus(t *testing.T) {
	defer func() { shared.ResetRunner() }()

	// Create a temp dir for postgres config
	tmpDir := t.TempDir()

	configPath := path.Join(tmpDir, "postgresql.conf")
	systemdDropin := path.Join(tmpDir, pgsql.BackupVolumeConfigName)

	type testCase struct {
		name           string
		archiveMode    string
		archiveCommand string
		mountPresent   bool
		mountConfig    string
		volumePresent  bool
		runnerOutput   map[string][]string
		expectError    string
	}

	cases := []testCase{
		{
			name:           "Enabled",
			archiveMode:    "on",
			archiveCommand: ArchiveCommand(),
			mountPresent:   true,
			volumePresent:  true,
			runnerOutput: map[string][]string{
				"psql": {"on", strings.Trim(ArchiveCommand(), "'")},
			},
		},
		{
			name:           "Disabled",
			archiveMode:    "off",
			archiveCommand: ArchiveCommand(),
			mountPresent:   true,
			volumePresent:  true,
			runnerOutput: map[string][]string{
				"psql": {"off", strings.Trim(ArchiveCommand(), "'")},
			},
			expectError: ErrArchiveModeOff.Error(),
		},
		{
			name:           "MisconfiguredCommand",
			archiveMode:    "on",
			archiveCommand: "/bin/false",
			mountPresent:   true,
			volumePresent:  true,
			runnerOutput: map[string][]string{
				"psql": {"on", "/bin/false"},
			},
			expectError: ErrArchiveCommandMisconfigured.Error(),
		},
		{
			name:           "MissingVolume",
			archiveMode:    "on",
			archiveCommand: ArchiveCommand(),
			mountPresent:   true,
			volumePresent:  false,
			runnerOutput: map[string][]string{
				"psql": {"on", strings.Trim(ArchiveCommand(), "'")},
			},
			expectError: ErrArchiveMountMisconfigured.Error(),
		},
		{
			name:           "MissingMount",
			archiveMode:    "on",
			archiveCommand: ArchiveCommand(),
			mountPresent:   false,
			volumePresent:  true,
			runnerOutput: map[string][]string{
				"psql": {"on", strings.Trim(ArchiveCommand(), "'")},
			},
			expectError: ErrArchiveMountMisconfigured.Error(),
		},
		{
			name:           "MisconfiguredMount",
			archiveMode:    "on",
			archiveCommand: ArchiveCommand(),
			mountPresent:   true,
			mountConfig:    "false",
			volumePresent:  true,
			runnerOutput: map[string][]string{
				"psql": {"on", strings.Trim(ArchiveCommand(), "'")},
			},
			expectError: ErrArchiveMountMisconfigured.Error(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock connection runner
			shared.SetRunner(func(command string, args ...string) types.Runner {
				// Handle psql commands
				if command == "podman" && len(args) > 0 {
					switch args[0] {
					case "exec":
						// Extract psql command
						// args: exec uyuni-db /usr/bin/psql -U postgres -tAc SHOW archive_mode;
						// or: exec uyuni-db /usr/bin/psql -U postgres -tAc SHOW archive_command;
						lastArg := args[len(args)-1]

						val := ""
						if strings.Contains(lastArg, "SHOW archive_mode;") {
							val = tc.archiveMode
						}
						if strings.Contains(lastArg, "SHOW archive_command;") {
							val = strings.Trim(tc.archiveCommand, "'")
						}
						return testutils.FakeRunnerGenerator(val, nil)(command, args...)
					}
				}
				return testutils.FakeRunnerGenerator("", nil)(command, args...)
			})

			// Mock host podman runner
			podman.SetRunner(func(command string, args ...string) types.Runner {
				if command == "podman" && len(args) > 0 {
					switch args[0] {
					case "volume":
						// podman volume inspect for GetVolumeMountPoint (called by ParsePostgresConfig)
						if args[1] == "inspect" {
							return testutils.FakeRunnerGenerator(tmpDir, nil)(command, args...)
						}
						// podman volume inspect for IsVolumeExists (called by CheckMount)
						if args[1] == "exists" {
							if tc.volumePresent {
								return testutils.FakeRunnerGenerator("", nil)(command, args...)
							}
							return testutils.FakeRunnerGenerator("", &exec.ExitError{})(command, args...)
						}
					case "inspect":
						// Handle podman inspect for GetPodName
						return testutils.FakeRunnerGenerator(podman.DBContainerName, nil)(command, args...)
					}
				}
				return testutils.FakeRunnerGenerator("", nil)(command, args...)
			})
			defer podman.ResetRunner()

			// Mock systemd
			originalSystemd := systemd
			defer func() { systemd = originalSystemd }()

			mockSystemd := &testutils.FakeSystemdDriver{
				ServiceProperties: map[string]map[string]string{
					podman.DBService: {
						podman.DropInPaths: "",
					},
				},
			}
			if tc.mountPresent {
				mockSystemd.ServiceProperties = map[string]map[string]string{
					podman.DBService: {
						podman.DropInPaths: "DropInPaths=" + systemdDropin,
					},
				}
			}
			systemd = podman.NewSystemdWithDriver(mockSystemd)

			// Prepare file config
			// checkStatusFile reads postgresql.conf
			configContent := []string{
				"archive_mode = " + tc.archiveMode,
				"archive_command = " + tc.archiveCommand,
			}
			testutils.WriteFile(t, configPath, strings.Join(configContent, "\n"))

			mountConfig := pgsql.BackupVolumeConfig()
			if tc.mountConfig != "" {
				mountConfig = tc.mountConfig
			}
			testutils.WriteFile(t, systemdDropin, mountConfig)

			// Run CheckStatus
			testutils.AssertError(t, tc.expectError, CheckStatus())
		})
	}
}
