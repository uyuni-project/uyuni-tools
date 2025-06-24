// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
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
	if err := os.Mkdir(serviceConfDir, 0750); err != nil {
		t.Fatalf("failed to create fake service configuration directory: %s", err)
	}

	servicesPath = testDir

	testutils.WriteFile(t, path.Join(serviceConfDir, "Service.conf"), currentFile)

	if err := CleanSystemdConfFile("uyuni-server"); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	actual := testutils.ReadFile(t, path.Join(serviceConfDir, "generated.conf"))
	testutils.AssertEquals(t, "invalid generated.conf file", generatedFile, actual)

	actual = testutils.ReadFile(t, path.Join(serviceConfDir, "custom.conf"))
	testutils.AssertEquals(t, "invalid custom.conf file", customFile, actual)

	if utils.FileExists(path.Join(serviceConfDir, "Service.conf")) {
		t.Error("the old Service.conf file is not removed")
	}
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
	if err := os.Mkdir(serviceConfDir, 0750); err != nil {
		t.Fatalf("failed to create fake service configuration directory: %s", err)
	}

	servicesPath = testDir

	testutils.WriteFile(t, path.Join(serviceConfDir, "generated.conf"), generatedFile)
	testutils.WriteFile(t, path.Join(serviceConfDir, "custom.conf"), customFile)

	if err := CleanSystemdConfFile("uyuni-server"); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	actual := testutils.ReadFile(t, path.Join(serviceConfDir, "generated.conf"))
	testutils.AssertEquals(t, "invalid generated.conf file", generatedFile, actual)

	actual = testutils.ReadFile(t, path.Join(serviceConfDir, "custom.conf"))
	testutils.AssertEquals(t, "invalid custom.conf file", customFile, actual)
}

func TestGetServiceProperty(t *testing.T) {
	newRunner = testutils.FakeRunnerGenerator("TestProperty=foo bar\n", nil)
	tested := SystemdImpl{}
	actual, err := tested.GetServiceProperty("myservice", "TestProperty")
	testutils.AssertTrue(t, "No error expected", err == nil)
	testutils.AssertEquals(t, "Wrong expected property", "foo bar", actual)
}

func TestGetServicePropertyError(t *testing.T) {
	newRunner = testutils.FakeRunnerGenerator("", errors.New("Test error"))
	tested := SystemdImpl{}
	actual, err := tested.GetServiceProperty("myservice", "TestProperty")
	testutils.AssertTrue(t, "Error message missing the root error message", strings.Contains(err.Error(), "Test error"))
	testutils.AssertTrue(t, "Unexpected error description",
		strings.Contains(err.Error(), "Failed to get the TestProperty property from myservice service"))
	testutils.AssertEquals(t, "Wrong expected property", "", actual)
}
