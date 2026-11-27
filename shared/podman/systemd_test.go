// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestCleanSystemdConfFile(t *testing.T) {
	currentFile := `[Service]
# Some comment
Environment=TZ=Europe/Berlin
Environment="PODMAN_EXTRA_ARGS="
Environment=UYUNI_IMAGE=path/to/image
`

	generatedFile := confHeader + `[Service]
Environment=UYUNI_IMAGE=path/to/image
`

	customFile := `[Service]
# Some comment
Environment=TZ=Europe/Berlin
Environment="PODMAN_EXTRA_ARGS="

`

	testDir := t.TempDir()

	serviceConfDir := path.Join(testDir, "uyuni-server.service.d")
	err := os.Mkdir(serviceConfDir, 0750)
	testutils.AssertNoError(t, "failed to create fake service configuration directory", err)

	servicesPath = testDir

	testutils.WriteFile(t, path.Join(serviceConfDir, "Service.conf"), currentFile)

	err = CleanSystemdConfFile("uyuni-server")
	testutils.AssertNoError(t, "unexpected error", err)

	actual := testutils.ReadFile(t, path.Join(serviceConfDir, GeneratedConf))
	testutils.AssertEquals(t, "invalid generated.conf file", generatedFile, actual)

	actual = testutils.ReadFile(t, path.Join(serviceConfDir, CustomConf))
	testutils.AssertEquals(t, "invalid custom.conf file", customFile, actual)

	testutils.AssertTrue(
		t, "the old Service.conf file is not removed", !utils.FileExists(path.Join(serviceConfDir, "Service.conf")),
	)
}

func TestCleanSystemdConfFileNoop(t *testing.T) {
	generatedFile := confHeader + `[Service]
Environment=UYUNI_IMAGE=path/to/image
`

	customFile := `[Service]
# Some comment
Environment=TZ=Europe/Berlin
Environment="PODMAN_EXTRA_ARGS="
`

	testDir := t.TempDir()

	serviceConfDir := path.Join(testDir, "uyuni-server.service.d")
	err := os.Mkdir(serviceConfDir, 0750)
	testutils.AssertNoError(t, "failed to create fake service configuration directory", err)

	servicesPath = testDir

	testutils.WriteFile(t, path.Join(serviceConfDir, GeneratedConf), generatedFile)
	testutils.WriteFile(t, path.Join(serviceConfDir, CustomConf), customFile)

	err = CleanSystemdConfFile("uyuni-server")
	testutils.AssertNoError(t, "unexpected error", err)

	actual := testutils.ReadFile(t, path.Join(serviceConfDir, GeneratedConf))
	testutils.AssertEquals(t, "invalid generated.conf file", generatedFile, actual)

	actual = testutils.ReadFile(t, path.Join(serviceConfDir, CustomConf))
	testutils.AssertEquals(t, "invalid custom.conf file", customFile, actual)
}

func TestGetServiceProperty(t *testing.T) {
	newRunner = testutils.FakeRunnerGenerator("TestProperty=foo bar\n", nil)
	tested := NewSystemd()
	actual, err := tested.GetServiceProperty("myservice", "TestProperty")
	testutils.AssertNoError(t, "No error expected", err)
	testutils.AssertEquals(t, "Wrong expected property", "foo bar", actual)
}

func TestGetServicePropertyError(t *testing.T) {
	newRunner = testutils.FakeRunnerGenerator("", errors.New("Test error"))
	tested := NewSystemd()
	actual, err := tested.GetServiceProperty("myservice", "TestProperty")
	testutils.AssertError(t, "Failed to get the TestProperty property from myservice service", err)
	testutils.AssertStringContains(t, "Error message missing the root error message", err.Error(), "Test error")
	testutils.AssertEquals(t, "Wrong expected property", "", actual)
}

func TestShowServiceLogs(t *testing.T) {
	var actualArgs []string
	SetRunner(func(command string, args ...string) types.Runner {
		actualArgs = append([]string{command}, args...)
		return testutils.FakeRunnerGenerator("some log line", nil)(command, args...)
	})
	t.Cleanup(ResetRunner)

	ShowServiceLogs("uyuni-server", time.Date(2026, 10, 26, 10, 0, 0, 0, time.UTC))

	expectedArgs := []string{
		"journalctl", "-u", "uyuni-server",
		"--since", "2026-10-26 10:00:00",
		"--no-pager", "-n", "100",
	}
	testutils.AssertEquals(t, "Wrong journalctl arguments", expectedArgs, actualArgs)
}

func TestShowServiceLogsFailure(t *testing.T) {
	SetRunner(testutils.FakeRunnerGenerator("", errors.New("no journal")))
	t.Cleanup(ResetRunner)

	// A failure to read the logs should not prevent reporting the service start failure.
	ShowServiceLogs("uyuni-server", time.Now())
}

func TestFormatServiceLogs(t *testing.T) {
	testutils.AssertEquals(t, "empty input should produce empty output", "", formatServiceLogs("myservice", ""))
	testutils.AssertEquals(t, "whitespace input should produce empty output", "",
		formatServiceLogs("myservice", "   \n\t  "))

	shortLog := "Started service successfully."
	expectedShort := "-------------------- [myservice] --------------------\n" +
		"  Started service successfully.\n" +
		"-----------------------------------------------------"
	testutils.AssertEquals(t, "short line format", expectedShort, formatServiceLogs("myservice", shortLog))

	longLine := strings.Repeat("x", maxJournalLineLength+50)
	formatted := formatServiceLogs("myservice", longLine)
	expectedTruncated := strings.Repeat("x", maxJournalLineLength-3) + "..."
	expectedLong := "-------------------- [myservice] --------------------\n" +
		"  " + expectedTruncated + "\n" +
		"-----------------------------------------------------"
	testutils.AssertEquals(t, "long line format with truncation", expectedLong, formatted)
}
