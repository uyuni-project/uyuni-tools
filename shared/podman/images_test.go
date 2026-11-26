// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestGetRpmImageName(t *testing.T) {
	data := []struct {
		name         string
		wantRPMImage string
		wantTag      string
		image        string
	}{
		{
			name:         "registry without tag",
			wantRPMImage: "suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			wantTag:      "latest",
			image:        "registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql",
		},
		{
			name:         "registry with latest tag",
			wantRPMImage: "suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			wantTag:      "latest",
			image:        "registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql:latest",
		},
		{
			name:         "http registry with version tag",
			wantRPMImage: "suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			wantTag:      "5.1.0",
			image:        "http://registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql:5.1.0",
		},
		{
			name:         "https registry with version tag",
			wantRPMImage: "suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			wantTag:      "5.1.0",
			image:        "https://registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql:5.1.0",
		},
		{
			name:         "docker registry with version tag",
			wantRPMImage: "suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			wantTag:      "5.1.0",
			image:        "docker://registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql:5.1.0",
		},
		{
			name:         "oci registry without tag",
			wantRPMImage: "suse-multi-linux-manager-5.1-x86_64-server-postgresql",
			wantTag:      "latest",
			image:        "oci://registry.suse.com/suse/multi-linux-manager/5.1/x86_64/server-postgresql",
		},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			rpmImageResult, tagResult := GetRpmImageName(tt.image)

			if tt.wantRPMImage != rpmImageResult {
				t.Errorf("Expected %s got %s when computing RPM for image %s", tt.wantRPMImage, rpmImageResult, tt.image)
			}
			if tt.wantTag != tagResult {
				t.Errorf("Expected %s got %s when computing RPM for image %s", tt.wantTag, tagResult, tt.image)
			}
		})
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

	data := []struct {
		name           string
		expectedResult string
		rpmImage       string
		tag            string
	}{
		{
			name:           "matching beta build tag",
			expectedResult: "/usr/share/suse-docker-images/native/multi-linux-manager-5.1-x86_64-server-5.1.0-beta1.x86_64-12.55.tar",
			rpmImage:       "multi-linux-manager-5.1-x86_64-server",
			tag:            "5.1.0-beta1.12.55",
		},
		{
			name:           "matching latest tag",
			expectedResult: "/usr/share/suse-docker-images/native/multi-linux-manager-5.1-x86_64-server-5.1.0-beta1.x86_64-12.55.tar",
			rpmImage:       "multi-linux-manager-5.1-x86_64-server",
			tag:            "latest",
		},
		{
			name:           "missing tag",
			expectedResult: "",
			rpmImage:       "multi-linux-manager-5.1-x86_64-server",
			tag:            "missing_tag",
		},
		{
			name:           "missing image and missing tag",
			expectedResult: "",
			rpmImage:       "missing_image",
			tag:            "missing_tag",
		},
		{
			name:           "missing image with latest tag",
			expectedResult: "",
			rpmImage:       "missing_image",
			tag:            "latest",
		},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			testResult, err := BuildRpmImagePath(jsonData, tt.rpmImage, tt.tag)

			if err != nil && tt.expectedResult != testResult {
				t.Errorf(
					"Expected %s got %s when computing RPM for image %s with tag %s",
					tt.expectedResult, testResult, tt.rpmImage, tt.tag,
				)
			}
		})
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
			res, err := PrepareImage("", tc.image, tc.pullPolicy, tc.pullEnabled)

			if tc.expectError {
				testutils.AssertError(t, fmt.Sprintf(tc.expectedMsg, tc.image), err)
				return
			}

			assertSuccess(t, res, tc.expectedImage, err)
		})
	}
}

func assertSuccess(t *testing.T, gotResult string, expectedResult string, err error) {
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedResult != "" {
		if gotResult != expectedResult {
			t.Errorf("Expected %q, got %q", expectedResult, gotResult)
		}
		return
	}

	if gotResult == "" {
		t.Errorf("Expected a valid image path, got empty string")
	}
}
