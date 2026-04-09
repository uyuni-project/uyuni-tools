// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"errors"
	"os"
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
	tmpDir, err := os.MkdirTemp("", "status_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

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
		expectError    error
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
			expectError: nil,
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
			expectError: ErrArchiveModeOff,
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
			expectError: ErrArchiveCommandMisconfigured,
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
			expectError: ErrArchiveMountMisconfigured,
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
			expectError: ErrArchiveMountMisconfigured,
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
			expectError: ErrArchiveMountMisconfigured,
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
			if err := os.WriteFile(configPath, []byte(strings.Join(configContent, "\n")), 0644); err != nil {
				t.Fatal(err)
			}

			mountConfig := pgsql.BackupVolumeConfig()
			if tc.mountConfig != "" {
				mountConfig = tc.mountConfig
			}
			if err := os.WriteFile(systemdDropin, []byte(mountConfig), 0600); err != nil {
				t.Fatal(err)
			}

			// Run CheckStatus
			err = CheckStatus()

			if tc.expectError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tc.expectError)
				} else if !errors.Is(err, tc.expectError) {
					t.Errorf("Expected error %v, got %v", tc.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected nil error, got %v", err)
				}
			}
		})
	}
}
