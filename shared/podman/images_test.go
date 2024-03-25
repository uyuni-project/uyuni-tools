// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"testing"
)

func TestGetRpmImageName(t *testing.T) {
	data := [][]string{
		{"suse-manager-5.0-x86_64-proxy-httpd", "latest", "registry.suse.com/suse/manager/5.0/x86_64/proxy-httpd"},
		{"suse-manager-5.0-x86_64-proxy-httpd", "latest", "registry.suse.com/suse/manager/5.0/x86_64/proxy-httpd:latest"},
		{"suse-manager-5.0-x86_64-proxy-httpd", "beta1", "registry.suse.com/suse/manager/5.0/x86_64/proxy-httpd:beta1"},
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
			"name": "suse-manager-5.0-x86_64-proxy-tftpd",
			"tags": ["latest", "5.0.0-beta1", "5.0.0-beta1.59.128"],
			"file": "suse-manager-5.0-x86_64-proxy-tftpd-latest.x86_64-59.128.tar"
		}
	}`)

	data := [][]string{
		{"/usr/share/suse-docker-images/native/suse-manager-5.0-x86_64-proxy-tftpd-latest.x86_64-59.128.tar", "suse-manager-5.0-x86_64-proxy-httpd", "latest"},
		{"/usr/share/suse-docker-images/native/suse-manager-5.0-x86_64-proxy-tftpd-latest.x86_64-59.128.tar", "suse-manager-5.0-x86_64-proxy-httpd", "5.0.0-beta1.59.128"},
		{"", "suse-manager-5.0-x86_64-proxy-httpd", "missing_tag"},
		{"", "missing_image", "missing_tag"},
		{"", "missing_image", "latest"},
	}

	for i, testCase := range data {
		expectedResult := testCase[0]
		rpmImage := testCase[1]
		tag := testCase[2]

		testResult, err := BuildRpmImagePath(jsonData, rpmImage, tag)

		if err != nil && expectedResult != testResult {
			t.Errorf("Testcase %d: Expected %s got %s when computing RPM for image %s with tag %s", i, expectedResult, testResult, rpmImage, tag)
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
