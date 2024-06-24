// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestGetContainerImage(t *testing.T) {
	data := [][]string{
		// Expectect image, value of --registry-server, value of --registry-path, value of --tag, value of --http-image, value of --http-tag
		{"registry/default/image/proxy-httpd:tag", "registry", "", "tag", "registry/default/image/proxy-httpd", ""},
		{"myregistry.example.com/proxy-httpd:tag", "myregistry.example.com", "", "tag", "default/image/proxy-httpd", ""},
		{"default/image/proxy-httpd:mytag", "", "default/image", "tag", "default/image/proxy-httpd", "mytag"},
	}

	for i, testCase := range data {
		proxyFlags := ProxyImageFlags{
			Tag:      testCase[3],
			Registry: types.RegistryFlags{Server: testCase[1], Path: testCase[2]},
			Httpd: types.ImageFlags{
				Name: testCase[4],
				Tag:  testCase[5],
			},
		}
		imageName := proxyFlags.GetContainerImage("httpd")
		if imageName != testCase[0] {
			t.Errorf("Testcase %d: Image name %s does not match expected %s", i, imageName, testCase[0])
		}
	}
}
