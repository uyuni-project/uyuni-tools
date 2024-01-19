// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "testing"

func TestComputeImage(t *testing.T) {
	data := [][]string{
		{"registry:5000/path/to/image:foo", "bar", "registry:5000/path/to/image:foo"},
		{"registry:5000/path/to/image", "bar", "registry:5000/path/to/image:bar"},
		{"registry/path/to/image:foo", "bar", "registry/path/to/image:foo"},
		{"registry/path/to/image", "bar", "registry/path/to/image:bar"},
	}

	for _, testCase := range data {
		image := testCase[0]
		tag := testCase[1]
		actual, err := ComputeImage(image, tag)
		if err != nil {
			t.Errorf("Unexpected error while computing image with %s, %s: %s", image, tag, err)
		}
		if actual != testCase[2] {
			t.Errorf("Expected %s got %s when computing image with %s, %s", testCase[2], actual, image, tag)
		}
	}
}

func TestComputeImageError(t *testing.T) {
	_, err := ComputeImage("registry:path/to/image:tag:tag", "bar")
	if err == nil {
		t.Error("Expected error, got none")
	}
}
