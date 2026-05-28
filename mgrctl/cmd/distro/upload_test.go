// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/mocks"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

const user = "testUser"
const password = "testPwd"
const server = "testServer"

var connectionDetails = &api.ConnectionDetails{User: user, Password: password, Server: server}

type distroUploadRequest struct {
	Filename string
	Distro   []byte
}

func readDistroUploadRequest(t *testing.T, req *http.Request) distroUploadRequest {
	t.Helper()

	reader, err := req.MultipartReader()
	if err != nil {
		t.Fatalf("Failed to create multipart reader: %v", err)
	}

	var data distroUploadRequest
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("Failed to read multipart part: %v", err)
		}

		partData, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("Failed to read multipart data: %v", err)
		}

		switch part.FormName() {
		case "filename":
			data.Filename = string(partData)
		case "distro":
			data.Distro = partData
			testutils.AssertEquals(t, "The form file name is not properly passed", data.Filename, part.FileName())
		}
	}

	return data
}

func TestDistroUpload(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		distro        []byte
		statusCode    int
		body          string
		expectedError string
	}{
		{
			name:          "Test uploading a distro ISO",
			filename:      "test.iso",
			distro:        []byte("test distro ISO content"),
			statusCode:    200,
			body:          `{"success":true,"result":1}`,
			expectedError: "",
		},
		{
			name:          "Test server returns error status",
			filename:      "test.iso",
			distro:        []byte("test distro ISO content"),
			statusCode:    500,
			body:          ``,
			expectedError: "error uploading distro: 500:",
		},
		{
			name:          "Test server reports upload failure",
			filename:      "test.iso",
			distro:        []byte("test distro ISO content"),
			statusCode:    200,
			body:          `{"success":false,"message":"invalid distro format"}`,
			expectedError: "failed to upload distro: invalid distro format",
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
					testutils.AssertEquals(t, "Wrong URL called", req.URL.Path, "/rhn/manager/api/admin/distro/uploadDistro")
					testutils.AssertTrue(t, "Wrong content type", strings.HasPrefix(req.Header.Get("Content-Type"),
						"multipart/form-data; boundary="))

					data := readDistroUploadRequest(t, req)
					testutils.AssertEquals(t, "The filename is not properly passed", tt.filename, data.Filename)
					testutils.AssertEquals(t, "The distro is not properly passed", tt.distro, data.Distro)

					return testutils.GetResponse(tt.statusCode, tt.body)
				},
			}

			errorMessage := ""
			if err := distroUpload(client, tt.filename, tt.distro); err != nil {
				errorMessage = err.Error()
			}
			testutils.AssertStringContains(t, "Unexpected error message", errorMessage, tt.expectedError)
		})
	}
}

func TestGetFilenameFromSource(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "Test local file",
			source:   "/tmp/test.iso",
			expected: "test.iso",
		},
		{
			name:     "Test URL",
			source:   "https://example.com/images/test.iso",
			expected: "test.iso",
		},
		{
			name:     "Test URL without filename",
			source:   "https://example.com/",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getFilenameFromSource(tt.source)
			testutils.AssertEquals(t, "Unexpected filename", tt.expected, actual)
		})
	}
}
