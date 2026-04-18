// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package knownhost

import (
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

// Test removeKnownHost function.
func TestRemoveKnownHost(t *testing.T) {
	tests := []struct {
		name          string
		hostname      string
		port          string
		statusCode    int
		body          string
		expectedError string
	}{
		{
			name:          "Test removing normal hostname and port",
			hostname:      "client.uyuni.lan",
			port:          "22",
			statusCode:    200,
			body:          `{"success":true,"result":1}`,
			expectedError: "",
		},
		{
			name:          "Test removing escaped hostname and port",
			hostname:      "escaped?client.uyuni.lan",
			port:          "42",
			statusCode:    200,
			body:          `{"success":true,"result":1}`,
			expectedError: "",
		},
		{
			name:          "Test removing an invalid port value",
			hostname:      "myotherclient.uyuni.lan",
			port:          "abc",
			statusCode:    400,
			body:          `{"success":false,"message":"mocked error"}`,
			expectedError: "mocked error",
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
					testutils.AssertEquals(t, "Wrong URL called", req.URL.Path, "/rhn/manager/api/admin/ssh/removeKnownHost")

					query := req.URL.Query()
					testutils.AssertEquals(t, "The hostname is not properly passed", tt.hostname, query.Get("hostname"))
					testutils.AssertEquals(t, "The port is not properly passed", tt.port, query.Get("port"))

					return testutils.GetResponse(tt.statusCode, tt.body)
				},
			}

			errorMessage := ""
			if err := removeKnownHost(client, tt.hostname, tt.port); err != nil {
				errorMessage = err.Error()
			}
			testutils.AssertStringContains(t, "Unexpected error message", errorMessage, tt.expectedError)
		})
	}
}
