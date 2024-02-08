// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "testing"

func TestComputeImage(t *testing.T) {
	data := [][]string{
		{"registry:5000/path/to/image:foo", "registry:5000/path/to/image:foo", "bar"},
		{"registry:5000/path/to/image:bar", "registry:5000/path/to/image", "bar"},
		{"registry/path/to/image:foo", "registry/path/to/image:foo", "bar"},
		{"registry/path/to/image:bar", "registry/path/to/image", "bar"},
		{"registry:5000/path/to/image-migration-14-16:foo", "registry:5000/path/to/image:foo", "bar", "-migration-14-16"},
		{"registry:5000/path/to/image-migration-14-16:bar", "registry:5000/path/to/image", "bar", "-migration-14-16"},
		{"registry/path/to/image-migration-14-16:foo", "registry/path/to/image:foo", "bar", "-migration-14-16"},
		{"registry/path/to/image-migration-14-16:bar", "registry/path/to/image", "bar", "-migration-14-16"},
	}

	for i, testCase := range data {
		result := testCase[0]
		image := testCase[1]
		tag := testCase[2]
		appendToImage := testCase[3:]

		actual, err := ComputeImage(image, tag, appendToImage...)

		if err != nil {
			t.Errorf("Testcase %d: Unexpected error while computing image with %s, %s, %s: %s", i, image, tag, appendToImage, err)
		}
		if actual != result {
			t.Errorf("Testcase %d: Expected %s got %s when computing image with %s, %s, %s", i, result, actual, image, tag, appendToImage)
		}
	}
}

func TestComputeImageError(t *testing.T) {
	_, err := ComputeImage("registry:path/to/image:tag:tag", "bar")
	if err == nil {
		t.Error("Expected error, got none")
	}
}
