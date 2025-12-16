// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetRpmImageName(t *testing.T) {
	data := [][]string{
		{
			"suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			"latest",
			"registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql"},
		{
			"suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			"latest",
			"registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql:latest"},
		{
			"suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			"5.1.0",
			"http://registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql:5.1.0",
		},
		{
			"suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			"5.1.0",
			"https://registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql:5.1.0",
		},
		{
			"suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			"5.1.0",
			"docker://registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql:5.1.0",
		},
		{
			"suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			"latest",
			"oci://registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql"},
	}

	for i, testCase := range data {
		rpmImage := testCase[0]
		tag := testCase[1]
		image := testCase[2]

		rpmImageResult, tagResult := GetRpmImageName(image)

		if rpmImage != rpmImageResult {
			t.Errorf("Testcase %d: Expected %s got %s when computing RPM for image %s", i, rpmImage, rpmImageResult, image)
		}
		if tag != tagResult {
			t.Errorf("Testcase %d: Expected %s got %s when computing RPM for image %s", i, tag, tagResult, image)
		}
	}
}

func TestMatchingMetadata(t *testing.T) {
	jsonData := []byte(`{
  "image": {
    "name": "multi-linux-manager-5.1-x86_64-server",
    "tags": [   "5.1.0-beta1",   "5.1.0-beta1.12.55",   "latest" ] ,
    "file": "multi-linux-manager-5.1-x86_64-server-5.1.0-beta1.x86_64-12.55.tar"
  }
}
`)

	data := [][]string{
		{
			"/usr/share/suse-docker-images/native/multi-linux-manager-5.1-x86_64-server-5.1.0-beta1.x86_64-12.55.tar",
			"multi-linux-manager-5.1-x86_64-server",
			"5.1.0-beta1.12.55",
		},
		{
			"/usr/share/suse-docker-images/native/multi-linux-manager-5.1-x86_64-server-5.1.0-beta1.x86_64-12.55.tar",
			"multi-linux-manager-5.1-x86_64-server",
			"latest",
		},
		{"", "multi-linux-manager-5.1-x86_64-server", "missing_tag"},
		{"", "missing_image", "missing_tag"},
		{"", "missing_image", "latest"},
	}

	for i, testCase := range data {
		expectedResult := testCase[0]
		rpmImage := testCase[1]
		tag := testCase[2]

		testResult, err := BuildRpmImagePath(jsonData, rpmImage, tag)

		if err != nil && expectedResult != testResult {
			t.Errorf(
				"Testcase %d: Expected %s got %s when computing RPM for image %s with tag %s",
				i, expectedResult, testResult, rpmImage, tag,
			)
		}
	}

	jsonDataInvalidWithTypo := []byte(`{
		"image: {
			"name": "suse-manager-5.0-x86_64-proxy-tftpd",
			"tags": ["latest", "5.0.0-beta1", "5.0.0-beta1.59.128"],
			"file": "suse-manager-5.0-x86_64-proxy-tftpd-latest.x86_64-59.128.tar"
		}
	}`)

	_, err := BuildRpmImagePath(jsonDataInvalidWithTypo, "", "")
	if err == nil {
		t.Error("typo in json: this should fail")
	}
}

func TestPrepareImage(t *testing.T) {
	tempDir := t.TempDir()

	origRpmDir := rpmImageDir
	rpmImageDir = filepath.Join(tempDir, "rpms")
	_ = os.MkdirAll(rpmImageDir, 0755)
	defer func() { rpmImageDir = origRpmDir }()

	tests := []struct {
		name          string
		image         string
		pullPolicy    string
		pullEnabled   bool
		expectedImage string
		expectError   bool
		expectedMsg   string
	}{
		//testing just the error
		{
			name:          "test without tag",
			image:         "registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql",
			pullPolicy:    "IfNotPresent",
			pullEnabled:   true,
			expectedImage: "",
			expectError:   true,
			expectedMsg:   "Cannot prepare image %s because tag is missing",
		},
		{
			name:          "test without registry",
			image:         "suse/multi-linux-manager/5.1/x86_64/server-postgresql:5.1.0",
			pullPolicy:    "IfNotPresent",
			pullEnabled:   true,
			expectedImage: "",
			expectError:   true,
			expectedMsg:   "Cannot prepare image %s because registry is missing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Run the function under test
			// The PATH is already modified by setupFakeBin, so it calls our script
			res, err := PrepareImage("", tc.image, tc.pullPolicy, tc.pullEnabled)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got success")
				}
				if tc.expectedMsg != "" {
					if !strings.Contains(err.Error(), tc.expectedMsg) {
						t.Errorf("Expected error message to contain %q, got %q", tc.expectedMsg, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tc.name == "Image Missing, RPM Available" {
					if res == "" {
						t.Errorf("Expected RPM path, got empty")
					}
				} else if res != tc.expectedImage {
					t.Errorf("Expected %q, got %q", tc.expectedImage, res)
				}
			}
		})
	}
}
