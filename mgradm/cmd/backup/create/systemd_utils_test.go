// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package create

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestFindService(t *testing.T) {
	oldSystemd := systemd
	defer func() { systemd = oldSystemd }()

	tests := []struct {
		name            string
		serviceName     string
		hasServiceMock  func(string) bool
		expectedService string
		expectedSkip    bool
	}{
		{
			name:        "Service directly exists",
			serviceName: "uyuni-server",
			hasServiceMock: func(name string) bool {
				return name == "uyuni-server"
			},
			expectedService: "uyuni-server",
			expectedSkip:    false,
		},
		{
			name:        "Service exists as template service",
			serviceName: "uyuni-server",
			hasServiceMock: func(name string) bool {
				return name == "uyuni-server@"
			},
			expectedService: "uyuni-server@",
			expectedSkip:    false,
		},
		{
			name:        "Service does not exist",
			serviceName: "uyuni-server",
			hasServiceMock: func(_ string) bool {
				return false
			},
			expectedService: "uyuni-server@",
			expectedSkip:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installedServices := []string{}
			if tt.hasServiceMock(tt.serviceName) {
				installedServices = append(installedServices, tt.serviceName)
			}
			if tt.hasServiceMock(tt.serviceName + "@") {
				installedServices = append(installedServices, tt.serviceName+"@")
			}

			systemd = podman.NewSystemdWithDriver(&testutils.FakeSystemdDriver{Installed: installedServices})
			service, skip := findService(tt.serviceName)
			testutils.AssertEquals(t, "service name mismatch", tt.expectedService, service)
			testutils.AssertEquals(t, "skip mismatch", tt.expectedSkip, skip)
		})
	}
}

func TestGatherSystemdItems(t *testing.T) {
	oldSystemd := systemd
	oldUyuniServices := utils.UyuniServices
	defer func() {
		systemd = oldSystemd
		utils.UyuniServices = oldUyuniServices
	}()

	// We'll use a single custom service for testing
	utils.UyuniServices = []types.UyuniService{
		{Name: "test-service"},
	}

	tempDir := t.TempDir()

	// Service is skipped (does not exist)
	systemd = podman.NewSystemdWithDriver(&testutils.FakeSystemdDriver{})
	items := gatherSystemdItems()
	testutils.AssertEquals(t, "items length should be 0", 0, len(items))

	// Service exists, but FragmentPath fails
	systemd = podman.NewSystemdWithDriver(&testutils.FakeSystemdDriver{
		Installed: []string{"test-service"},
	})
	items = gatherSystemdItems()
	testutils.AssertEquals(t, "items length should be 0 on FragmentPath error", 0, len(items))

	// Service exists, FragmentPath succeeds, but DropInPaths fails, and no env file
	servicePath := filepath.Join(tempDir, "test-service.service")
	systemd = podman.NewSystemdWithDriver(&testutils.FakeSystemdDriver{
		Installed: []string{"test-service"},
		ServiceProperties: map[string]map[string]string{
			"test-service": {
				podman.FragmentPath: servicePath,
			},
		},
	})
	items = gatherSystemdItems()
	expectedItems := []string{servicePath}
	testutils.AssertEquals(t, "items should only contain service path", expectedItems, items)

	// Service exists, FragmentPath succeeds, DropInPaths is present but empty, and no env file
	systemd = podman.NewSystemdWithDriver(&testutils.FakeSystemdDriver{
		Installed: []string{"test-service"},
		ServiceProperties: map[string]map[string]string{
			"test-service": {
				podman.FragmentPath: servicePath,
				podman.DropInPaths:  "",
			},
		},
	})
	items = gatherSystemdItems()
	expectedItems = []string{servicePath}
	testutils.AssertEquals(t, "items should only contain service path when DropInPaths is empty", expectedItems, items)

	// Service exists, FragmentPath succeeds, DropInPaths succeeds, and env file exists
	serviceDir := servicePath + ".d"
	err := os.MkdirAll(serviceDir, 0755)
	testutils.AssertNoError(t, "failed to create service directory", err)

	envFile := filepath.Join(serviceDir, podman.ServerEnvironmentFile)
	testutils.WriteFile(t, envFile, "some env")

	dropIn1 := filepath.Join(serviceDir, "10-test.conf")
	dropIn2 := filepath.Join(serviceDir, "20-test.conf")
	// Test with a drop-in file that has spaces in its name. Systemd adds quotes around such paths
	dropIn3 := filepath.Join(serviceDir, "30 with space test.conf")
	dropInPaths := dropIn1 + " " + dropIn2 + " " + "\"" + dropIn3 + "\""

	systemd = podman.NewSystemdWithDriver(&testutils.FakeSystemdDriver{
		Installed: []string{"test-service"},
		ServiceProperties: map[string]map[string]string{
			"test-service": {
				podman.FragmentPath: servicePath,
				podman.DropInPaths:  dropInPaths,
			},
		},
	})

	items = gatherSystemdItems()
	expectedItems = []string{
		servicePath,
		serviceDir,
		dropIn1,
		dropIn2,
		dropIn3,
		envFile,
	}
	testutils.AssertEquals(t, "items mismatch", expectedItems, items)
}

func TestExportSystemdConfiguration(t *testing.T) {
	oldSystemd := systemd
	oldUyuniServices := utils.UyuniServices
	defer func() {
		systemd = oldSystemd
		utils.UyuniServices = oldUyuniServices
	}()

	tempDir := t.TempDir()

	serviceFile := filepath.Join(tempDir, "fake-service.service")
	serviceContent := "[Service]\nExecStart=/usr/bin/fake"
	testutils.WriteFile(t, serviceFile, serviceContent)

	serviceDir := serviceFile + ".d"
	err := os.MkdirAll(serviceDir, 0755)
	testutils.AssertNoError(t, "failed to create service directory", err)

	dropInFile := filepath.Join(serviceDir, "custom.conf")
	dropInContent := "[Service]\nEnvironment=FOO=bar"
	testutils.WriteFile(t, dropInFile, dropInContent)

	dropInWhiteSpaceFile := filepath.Join(serviceDir, "with space test.conf")
	dropInWhiteSpaceContent := "[Service]\nEnvironment=BAR=foo"
	testutils.WriteFile(t, dropInWhiteSpaceFile, dropInWhiteSpaceContent)

	dropInPaths := dropInFile + " " + "\"" + dropInWhiteSpaceFile + "\""

	utils.UyuniServices = []types.UyuniService{
		{Name: "fake-service"},
	}

	envFileContent := "SOME_ENV=1"
	envFile := filepath.Join(serviceDir, podman.ServerEnvironmentFile)
	testutils.WriteFile(t, envFile, envFileContent)

	systemd = podman.NewSystemdWithDriver(&testutils.FakeSystemdDriver{
		Installed: []string{"fake-service"},
		ServiceProperties: map[string]map[string]string{
			"fake-service": {
				podman.FragmentPath: serviceFile,
				podman.DropInPaths:  dropInPaths,
			},
		},
	})

	outputDir := t.TempDir()

	// dryRun = true
	err = exportSystemdConfiguration(outputDir, true)
	testutils.AssertNoError(t, "exportSystemdConfiguration with dryRun=true failed", err)

	backupFilePath := filepath.Join(outputDir, "systemdBackup.tar")
	testutils.AssertTrue(t, "backup file should not have been created on dryRun=true", !utils.FileExists(backupFilePath))

	// dryRun = false
	err = exportSystemdConfiguration(outputDir, false)
	testutils.AssertNoError(t, "exportSystemdConfiguration with dryRun=false failed", err)

	testutils.AssertTrue(t, "backup file was not created on dryRun=false", utils.FileExists(backupFilePath))

	tarFile, err := os.Open(backupFilePath)
	testutils.AssertNoError(t, "Open backup file failed", err)
	defer tarFile.Close()

	tr := tar.NewReader(tarFile)
	foundFiles := make(map[string]string)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		testutils.AssertNoError(t, "tar next failed", err)

		if header.FileInfo().IsDir() {
			foundFiles[header.Name] = "dir"
		} else {
			content, err := io.ReadAll(tr)
			testutils.AssertNoError(t, "failed to read tar entry", err)
			foundFiles[header.Name] = string(content)
		}
	}

	expectedServiceDirName := serviceDir + "/"
	testutils.AssertEquals(t, "service file content", serviceContent, foundFiles[serviceFile])
	testutils.AssertEquals(t, "service dir entry type", "dir", foundFiles[expectedServiceDirName])
	testutils.AssertEquals(t, "drop-in file content", dropInContent, foundFiles[dropInFile])
	testutils.AssertEquals(t, "drop-in white space content", dropInWhiteSpaceContent, foundFiles[dropInWhiteSpaceFile])
	testutils.AssertEquals(t, "env file content", envFileContent, foundFiles[envFile])
}
