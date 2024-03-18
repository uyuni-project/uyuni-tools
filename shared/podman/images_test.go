// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"testing"
)

func TestCalculateRpmImagePath(t *testing.T) {
	data := [][]string{
		{"/usr/share/suse-docker-images/native/suse-manager-5.0-x86_64-proxy-httpd", "registry.suse.com/suse/manager/5.0/x86_64/proxy-httpd:latest"},
		{"/usr/share/suse-docker-images/native/suse-manager-5.0-x86_64-proxy-httpd", "registry.suse.com/suse/manager/5.0/x86_64/proxy-httpd:beta1"},
		{"/usr/share/suse-docker-images/native/suse-manager-5.0-x86_64-proxy-httpd", "registry.suse.com/suse/manager/5.0/x86_64/proxy-httpd"},
	}

	for i, testCase := range data {
		result := testCase[0]
		image := testCase[1]

		ret := calculateRpmImagePath(image)

		if ret != result {
			t.Errorf("Testcase %d: Expected %s got %s when computing RPM for image %s", i, result, ret, image)
		}
	}
}
