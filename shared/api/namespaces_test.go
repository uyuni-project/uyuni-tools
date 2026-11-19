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

const sshNamespace = "admin.ssh"
const removeKnownHostEndpoint = "removeKnownHost"

func getAPICallListResponse(namespace string, endpoint string) string {
	return `{
		"success":true,
		"result":{
			"` + namespace + `": {
				"endpointid": {
					"name": "` + endpoint + `",
					"parameters": [],
					"exceptions": [],
					"return": "int"
				}
			}
		}
	}`
}

func TestValidateEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectedErr string
	}{
		{
			name:        "Matching path only",
			path:        "admin/ssh/removeKnownHost",
			expectedErr: "",
		},
		{
			name:        "Matching path with parameters",
			path:        "admin/ssh/removeKnownHost?key=val",
			expectedErr: "",
		},
		{
			name:        "Non matching path",
			path:        "foo/bar",
			expectedErr: "This function is currently not supported by the server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := Init(&ConnectionDetails{Server: "testServer", Cookie: "testCookie"})
			if err != nil {
				t.FailNow()
			}
			client.Client = &mocks.MockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				switch req.URL.Path {
				case "/rhn/manager/api/api/getApiCallList":
					return testutils.GetResponse(200, getAPICallListResponse(sshNamespace, removeKnownHostEndpoint))
				default:
					t.Errorf("Unexpected API path %s", req.URL.Path)
					return testutils.GetResponse(404, `{}`)
				}
			}}

			err = client.validateEndpoint(tt.path)
			if tt.expectedErr == "" {
				testutils.AssertNoError(t, "Shouldn't have failed", err)
			} else {
				testutils.AssertError(t, tt.expectedErr, err)
			}
		})
	}
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
			apiPath:      "/rhn/manager/api/admin/ssh/removeKnownHost",
			targetCalled: true,
		},
		{
			name:          "Test unsupported endpoint",
			path:          "admin/ssh/invalidCall",
			apiPath:       "/rhn/manager/api/admin/ssh/invalidCall",
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
				case "/rhn/manager/api/api/getApiCallList":
					return testutils.GetResponse(200, getAPICallListResponse(sshNamespace, removeKnownHostEndpoint))
				case tt.apiPath:
					targetCalled = true
					return testutils.GetResponse(200, `{"success":true,"result":true}`)
				default:
					t.Errorf("Unexpected API path %s", req.URL.Path)
					return testutils.GetResponse(404, `{}`)
				}
			}}

			errorMessage := ""
			if _, err := client.PostChecked(tt.path, nil); err != nil {
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
			apiPath:      "/rhn/manager/api/org/getDetails",
			targetCalled: true,
		},
		{
			name:          "Test unsupported endpoint",
			path:          "org/getUnknownDetails?name=admin",
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
				case "/rhn/manager/api/api/getApiCallList":
					return testutils.GetResponse(200, getAPICallListResponse("org", "getDetails"))
				case tt.apiPath:
					targetCalled = true
					return testutils.GetResponse(200, `{"success":true,"result":true}`)
				default:
					t.Errorf("Unexpected API path %s", req.URL.Path)
					return testutils.GetResponse(404, `{}`)
				}
			}}

			errorMessage := ""
			if _, err := client.GetChecked(tt.path); err != nil {
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
