// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package gpg

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/mocks"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

const user = "testUser"
const password = "testPwd"
const server = "testServer"

var connectionDetails = &api.ConnectionDetails{User: user, Password: password, Server: server}

type gpgKeyUploadRequest struct {
	Key string `json:"gpgKey"`
}

// Test gpgKeyUpload function.
func TestGpgKeyUpload(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		statusCode    int
		body          string
		expectedError string
	}{
		{
			name:          "Test uploading a GPG key",
			key:           "-----BEGIN PGP PUBLIC KEY BLOCK-----\ntest\n-----END PGP PUBLIC KEY BLOCK-----",
			statusCode:    200,
			body:          `{"success":true,"result":1}`,
			expectedError: "",
		},
		{
			name:          "Test uploading a GPG key with special characters",
			key:           "-----BEGIN PGP PUBLIC KEY BLOCK-----\nkey+with/special=chars\n-----END PGP PUBLIC KEY BLOCK-----",
			statusCode:    200,
			body:          `{"success":true,"result":1}`,
			expectedError: "",
		},
		{
			name:          "Test server returns error status",
			key:           "-----BEGIN PGP PUBLIC KEY BLOCK-----\ntest\n-----END PGP PUBLIC KEY BLOCK-----",
			statusCode:    500,
			body:          ``,
			expectedError: "error uploading GPG key: 500:",
		},
		{
			name:          "Test server reports upload failure",
			key:           "-----BEGIN ABC PUBLIC KEY BLOCK-----\ntest\n-----END PGP PUBLIC KEY BLOCK-----",
			statusCode:    200,
			body:          `{"success":false,"message":"invalid key format"}`,
			expectedError: "failed to upload GPG key: invalid key format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := api.Init(connectionDetails)
			if err != nil {
				t.FailNow()
			}

			client.Client = &mocks.MockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					testutils.AssertEquals(t, "Wrong URL called", req.URL.Path, "/rhn/manager/api/admin/gpg/uploadGpgKey")
					testutils.AssertEquals(t, "Wrong content type", req.Header.Get("Content-Type"), "application/json; charset=utf-8")

					var data gpgKeyUploadRequest
					if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
						t.Fatalf("Failed to decode JSON body: %v", err)
					}
					testutils.AssertEquals(t, "The key is not properly passed", tt.key, data.Key)

					return testutils.GetResponse(tt.statusCode, tt.body)
				},
			}

			errorMessage := ""
			if err := gpgKeyUpload(client, tt.key); err != nil {
				errorMessage = err.Error()
			}
			testutils.AssertStringContains(t, "Unexpected error message", errorMessage, tt.expectedError)
		})
	}
}
