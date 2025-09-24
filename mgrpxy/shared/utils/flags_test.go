// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// TestGetContainerImage tests GetContainerImage method
// Covering different scenarios: defaults, empty values, and overriding values.
func TestGetContainerImage(t *testing.T) {
	tests := []struct {
		name           string
		proxyFlags     ProxyImageFlags
		expectedResult string
		description    string
	}{
		// Defaults and overiding values
		{
			name: "no image details",
			proxyFlags: ProxyImageFlags{
				Registry: types.Registry{
					Host: "default/image",
				},
				Tag: "tag",
				Httpd: types.ImageFlags{
					Name: "",
					Tag:  "",
				},
			},
			expectedResult: "default/image:tag",
		},
		{
			name: "httpd image with registry",
			proxyFlags: ProxyImageFlags{
				Registry: types.Registry{
					Host: "default",
				},
				Tag: "tag",
				Httpd: types.ImageFlags{
					Name: "image/proxy-httpd",
					Tag:  "mytag",
				},
			},
			expectedResult: "default/image/proxy-httpd:mytag",
		},
		{
			name: "httpd image name is appended to registry when it does not include registry",
			proxyFlags: ProxyImageFlags{
				Registry: types.Registry{
					Host: "default/extra/image",
				},
				Tag: "tag",
				Httpd: types.ImageFlags{
					Name: "default/image/proxy-httpd",
					Tag:  "mytag",
				},
			},
			expectedResult: "default/extra/image/default/image/proxy-httpd:mytag",
		},

		// domain usage
		{
			name: "custom full httpd registry image name",
			proxyFlags: ProxyImageFlags{
				Registry: types.Registry{
					Host: "registry.suse.com/suse/some/paths/",
				},
				Tag: "1.0.0",
				Httpd: types.ImageFlags{
					Name: "uyuni/proxy-httpd",
					Tag:  "2.0.0",
				},
			},
			expectedResult: "registry.suse.com/suse/some/paths/uyuni/proxy-httpd:2.0.0",
			// expectedResult: "registry.opensuse.org/uyuni/proxy-httpd:2.0.0", // this should be the expected result
		},
		{
			name: "httpd with path-only image name",
			proxyFlags: ProxyImageFlags{
				Registry: types.Registry{
					Host: "registry.suse.com/uyuni",
				},
				Tag: "1.0.0",
				Httpd: types.ImageFlags{
					Name: "path/to/proxy-httpd",
					Tag:  "",
				},
			},
			expectedResult: "registry.suse.com/uyuni/path/to/proxy-httpd:1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.proxyFlags.GetContainerImage("httpd")

			if actual != tt.expectedResult {
				t.Errorf("GetContainerImage('httpd') = %s, expected: %s", actual, tt.expectedResult)
			}
		})
	}
}
