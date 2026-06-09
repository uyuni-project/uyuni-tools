// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/api/mocks"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

const removeKnownHostNamespace = "api.admin.ssh.remove_known_host"

func namespaceListResponse(namespace string) string {
	return `{"success":true,"result":[{"namespace":"` + namespace + `","access_mode":"W"}]}`
}

func TestPostChecked(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		endpoint      string
		apiPath       string
		expectedError string
		targetCalled  bool
	}{
		{
			name:         "Test advertised endpoint",
			path:         "admin/ssh/removeKnownHost",
			endpoint:     removeKnownHostNamespace,
			apiPath:      "/rhn/manager/api/admin/ssh/removeKnownHost",
			targetCalled: true,
		},
		{
			name:          "Test unsupported endpoint",
			path:          "admin/ssh/removeUnknownHost",
			endpoint:      "api.admin.ssh.remove_unknown_host",
			apiPath:       "/rhn/manager/api/admin/ssh/removeUnknownHost",
			expectedError: unsupportedFunctionError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := Init(&ConnectionDetails{Server: "testServer", Cookie: "testCookie"})
			if err != nil {
				t.FailNow()
			}

			targetCalled := false
			client.Client = &mocks.MockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				switch req.URL.Path {
				case "/rhn/manager/api/access/listNamespaces":
					return testutils.GetResponse(200, namespaceListResponse(removeKnownHostNamespace))
				case tt.apiPath:
					targetCalled = true
					return testutils.GetResponse(200, `{"success":true,"result":true}`)
				default:
					t.Errorf("Unexpected API path %s", req.URL.Path)
					return testutils.GetResponse(404, `{}`)
				}
			}}

			errorMessage := ""
			if _, err := client.PostChecked(tt.path, tt.endpoint, nil); err != nil {
				errorMessage = err.Error()
			}

			testutils.AssertStringContains(t, "Unexpected error message", errorMessage, tt.expectedError)
			testutils.AssertEquals(t, "Unexpected target call", tt.targetCalled, targetCalled)
		})
	}
}

func TestGetChecked(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		endpoint      string
		apiPath       string
		expectedError string
		targetCalled  bool
	}{
		{
			name:         "Test advertised endpoint",
			path:         "org/getDetails?name=admin",
			endpoint:     "api.org.get_details",
			apiPath:      "/rhn/manager/api/org/getDetails",
			targetCalled: true,
		},
		{
			name:          "Test unsupported endpoint",
			path:          "org/getUnknownDetails?name=admin",
			endpoint:      "api.org.get_unknown_details",
			apiPath:       "/rhn/manager/api/org/getUnknownDetails",
			expectedError: unsupportedFunctionError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := Init(&ConnectionDetails{Server: "testServer", Cookie: "testCookie"})
			if err != nil {
				t.FailNow()
			}

			targetCalled := false
			client.Client = &mocks.MockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				switch req.URL.Path {
				case "/rhn/manager/api/access/listNamespaces":
					return testutils.GetResponse(200, namespaceListResponse("api.org.get_details"))
				case tt.apiPath:
					targetCalled = true
					return testutils.GetResponse(200, `{"success":true,"result":true}`)
				default:
					t.Errorf("Unexpected API path %s", req.URL.Path)
					return testutils.GetResponse(404, `{}`)
				}
			}}

			errorMessage := ""
			if _, err := client.GetChecked(tt.path, tt.endpoint); err != nil {
				errorMessage = err.Error()
			}

			testutils.AssertStringContains(t, "Unexpected error message", errorMessage, tt.expectedError)
			testutils.AssertEquals(t, "Unexpected target call", tt.targetCalled, targetCalled)
		})
	}
}

func TestRawPostDoesNotValidateEndpoint(t *testing.T) {
	client, err := Init(&ConnectionDetails{Server: "testServer", Cookie: "testCookie"})
	if err != nil {
		t.FailNow()
	}

	targetCalled := false
	client.Client = &mocks.MockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		testutils.AssertEquals(t, "Wrong URL called", req.URL.Path, "/rhn/manager/api/org/createFirst")
		targetCalled = true
		return testutils.GetResponse(200, `{"success":true,"result":true}`)
	}}

	if _, err := client.Post("org/createFirst", nil); err != nil {
		t.Fatalf("Expected raw POST request to bypass endpoint validation: %v", err)
	}
	testutils.AssertTrue(t, "Expected raw target endpoint to be called", targetCalled)
}
