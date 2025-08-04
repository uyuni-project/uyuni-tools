// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestCheckDirPermissions(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := checkPermissions(tempDir, 0005|0050|0500); err != nil {
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
	if err := os.Chmod(configFilePath, 0600); err != nil {
		t.Fatalf("Failed to change permissions for %s: %v", configFilePath, err)
	}
	if err := validateInstallYamlFiles(tempDir); err == nil {
		t.Errorf("Expected an error due to incorrect permissions on config.yaml, but got none")
	}

	// Restore the correct permissions for the next test run
	if err := os.Chmod(configFilePath, 0644); err != nil {
		t.Fatalf("Failed to restore permissions for %s: %v", configFilePath, err)
	}

	// Test: Missing file scenario, remove one file and expect an error
	os.Remove(path.Join(tempDir, "httpd.yaml"))
	if err := validateInstallYamlFiles(tempDir); err == nil {
		t.Errorf("Expected an error due to missing httpd.yaml, but got none")
	}
}

func TestGetSystemID(t *testing.T) {
	// event output
	systemid := `<?xml version=\"1.0\"?><params><param><value><struct><member><name>username</name>` +
		`<value><string>admin</string></value></member><member><name>os_release</name><value><string>6.1</string>` +
		`</value></member><member><name>operating_system</name><value><string>SL-Micro</string></value></member>` +
		`<member><name>architecture</name><value><string>x86_64-redhat-linux</string></value></member><member>` +
		`<name>system_id</name><value><string>ID-1000010001</string></value></member><member><name>type</name><value>` +
		`<string>REAL</string></value></member><member><name>fields</name><value><array><data><value>` +
		`<string>system_id</string></value><value><string>os_release</string></value><value><string>operating_system` +
		`</string></value><value><string>architecture</string></value><value><string>username</string></value><value>` +
		`<string>type</string></value></data></array></value></member><member><name>checksum</name><value>` +
		`<string>1aaa4427328cfd7fbd613693802e0920d9f1c1ea2b3d31a869ed1ac3fbfe4174</string></value></member></struct>` +
		`</value></param></params>`

	event := `suse/systemid/generated {"data": "` + systemid + `", "_stamp": "2025-08-04T12:04:29.403745"}`

	// create custom runners
	contextRunner = testutils.FakeContextRunnerGenerator(event, nil)
	newRunner = testutils.FakeRunnerGenerator("", nil)

	received, err := getSystemIDEvent()
	testutils.AssertNoError(t, "error during obtaining systemid", err)
	testutils.AssertEquals(t, "received event differs", []byte(event), received)

	receivedSystemid, err := parseSystemIDEvent(received)
	testutils.AssertNoError(t, "error during event decoding", err)
	// unquote raw string before comparing.
	unquotedSystemid, _ := strconv.Unquote(`"` + systemid + `"`)
	testutils.AssertEquals(t, "received systemid differs", unquotedSystemid, receivedSystemid)
}
