// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
	Name     string
	HasName  bool
	Distro   []byte
}

func readDistroUploadRequest(t *testing.T, req *http.Request) distroUploadRequest {
	t.Helper()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("Failed to read request body: %v", err)
	}
	testutils.AssertEquals(t, "Content-Length doesn't match the body size", req.ContentLength, int64(len(body)))
	req.Body = io.NopCloser(bytes.NewReader(body))

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
		case "name":
			data.Name = string(partData)
			data.HasName = true
		case "distro":
			data.Distro = partData
			testutils.AssertEquals(t, "The form file name is not properly passed", data.Filename, part.FileName())
		}
	}

	return data
}

func TestDistroUpload(t *testing.T) {
	tests := []struct {
		name           string
		filename       string
		distroName     string
		distro         []byte
		statusCode     int
		body           string
		expectedError  string
		expectedOutput []string
	}{
		{
			name:          "Test uploading a distro ISO",
			filename:      "test.iso",
			distro:        []byte("test distro ISO content"),
			statusCode:    200,
			body:          `{"success":true,"result":1}`,
			expectedError: "",
			expectedOutput: []string{
				"Distro successfully uploaded to /srv/www/distributions/test\n",
			},
		},
		{
			name:          "Test uploading a distro ISO with a name to a legacy server",
			filename:      "rhel-server-7.0-x86_64-dvd.iso",
			distroName:    "rhel-server-7.0",
			distro:        []byte("test distro ISO content"),
			statusCode:    200,
			body:          `{"success":true,"result":1}`,
			expectedError: "",
			expectedOutput: []string{
				"Distro successfully uploaded to /srv/www/distributions/rhel-server-7.0\n",
				"Declare the distribution in the Web UI if needed",
			},
		},
		{
			name:       "Test uploading a distro ISO registered by the server",
			filename:   "test.iso",
			distro:     []byte("test distro ISO content"),
			statusCode: 200,
			body: `{"success":true,"result":` +
				`{"path":"/srv/www/distributions/test","name":"test","registered":"test-label"}}`,
			expectedError: "",
			expectedOutput: []string{
				"Distro successfully uploaded to /srv/www/distributions/test\n",
				"Distribution test-label successfully registered",
			},
		},
		{
			name:          "Test uploading a distro ISO not registered by the server",
			filename:      "test.iso",
			distro:        []byte("test distro ISO content"),
			statusCode:    200,
			body:          `{"success":true,"result":{"path":"/srv/www/distributions/other","name":"other","registered":""}}`,
			expectedError: "",
			expectedOutput: []string{
				"Distro successfully uploaded to /srv/www/distributions/other\n",
				"not registered automatically",
			},
		},
		{
			name:          "Test uploading an empty distro ISO",
			filename:      "test.iso",
			distro:        []byte{},
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
		{
			name:          "Test server returns unexpected result",
			filename:      "test.iso",
			distro:        []byte("test distro ISO content"),
			statusCode:    200,
			body:          `{"success":true,"result":0}`,
			expectedError: "unable to upload distro, server returned an error",
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
					testutils.AssertEquals(t, "The name is not properly passed", tt.distroName, data.Name)
					testutils.AssertEquals(t, "The name field should only be sent if set",
						tt.distroName != "", data.HasName)
					testutils.AssertEquals(t, "The distro is not properly passed", tt.distro, data.Distro)

					return testutils.GetResponse(tt.statusCode, tt.body)
				},
			}

			errorMessage := ""
			size := int64(len(tt.distro))
			output := captureStdout(t, func() {
				if err := distroUpload(client, tt.filename, tt.distroName, bytes.NewReader(tt.distro), size); err != nil {
					errorMessage = err.Error()
				}
			})
			testutils.AssertStringContains(t, "Unexpected error message", errorMessage, tt.expectedError)
			for _, expected := range tt.expectedOutput {
				testutils.AssertStringContains(t, "Unexpected output", output, expected)
			}
		})
	}
}

func TestDistroUploadUnsupportedServer(t *testing.T) {
	client, err := api.Init(&api.ConnectionDetails{Server: server, Cookie: "testCookie"})
	if err != nil {
		t.FailNow()
	}

	client.Client = &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			testutils.AssertEquals(t, "Only the API call list should be requested",
				"/rhn/manager/api/api/getApiCallList", req.URL.Path)
			return testutils.GetResponse(200, `{"success":true,"result":{"admin.gpg":{"id":{"name":"uploadGpgKey"}}}}`)
		},
	}

	distro := &failingReader{t: t}
	err = distroUpload(client, "test.iso", "", distro, 10)
	testutils.AssertTrue(t, "Upload should fail", err != nil)
	testutils.AssertStringContains(t, "Unexpected error message", err.Error(),
		"This function is currently not supported by the server")
}

// failingReader fails the test if the distro gets read.
type failingReader struct {
	t *testing.T
}

func (r *failingReader) Read(_ []byte) (int, error) {
	r.t.Error("The distro should not be read")
	return 0, io.EOF
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
			name:     "Test URL with query",
			source:   "https://example.com/images/test.iso?token=abc",
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

func TestOpenDistroFile(t *testing.T) {
	content := []byte("test distro ISO content")
	source := filepath.Join(t.TempDir(), "test.iso")
	if err := os.WriteFile(source, content, 0600); err != nil {
		t.Fatalf("Failed to write test ISO: %v", err)
	}

	reader, size, filename, err := openDistro(source)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to read distro: %v", err)
	}
	testutils.AssertEquals(t, "Unexpected filename", "test.iso", filename)
	testutils.AssertEquals(t, "Unexpected size", int64(len(content)), size)
	testutils.AssertEquals(t, "Unexpected content", content, data)
}

func TestOpenDistroErrors(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name          string
		source        string
		expectedError string
	}{
		{
			name:          "Test missing local file",
			source:        filepath.Join(dir, "missing.iso"),
			expectedError: "failed to read distro ISO file",
		},
		{
			name:          "Test directory",
			source:        dir,
			expectedError: "is not a regular file",
		},
		{
			name:          "Test URL without filename",
			source:        "https://example.com/",
			expectedError: "unable to determine distro ISO filename",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := openDistro(tt.source)
			testutils.AssertTrue(t, "openDistro should fail", err != nil)
			testutils.AssertStringContains(t, "Unexpected error message", err.Error(), tt.expectedError)
		})
	}
}

func TestOpenDistroURL(t *testing.T) {
	content := []byte("test distro ISO content")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/sized/test.iso":
			_, _ = w.Write(content)
		case "/chunked/test.iso":
			// Flushing before the end forces a chunked response without Content-Length
			_, _ = w.Write(content[:4])
			w.(http.Flusher).Flush()
			_, _ = w.Write(content[4:])
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	t.Run("Test URL with known size", func(t *testing.T) {
		reader, size, filename, err := openDistro(srv.URL + "/sized/test.iso")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("Failed to read distro: %v", err)
		}
		testutils.AssertEquals(t, "Unexpected filename", "test.iso", filename)
		testutils.AssertEquals(t, "Unexpected size", int64(len(content)), size)
		testutils.AssertEquals(t, "Unexpected content", content, data)
	})

	t.Run("Test URL with unknown size", func(t *testing.T) {
		_, _, _, err := openDistro(srv.URL + "/chunked/test.iso")
		testutils.AssertTrue(t, "openDistro should fail", err != nil)
		testutils.AssertStringContains(t, "Unexpected error message", err.Error(), "did not report the ISO size")
	})

	t.Run("Test URL not found", func(t *testing.T) {
		_, _, _, err := openDistro(srv.URL + "/missing/test.iso")
		testutils.AssertTrue(t, "openDistro should fail", err != nil)
		testutils.AssertStringContains(t, "Unexpected error message", err.Error(), "bad status 404")
	})
}

// captureStdout returns what the function prints on the standard output.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	stdout := os.Stdout
	os.Stdout = writer
	defer func() { os.Stdout = stdout }()

	output := make(chan string)
	go func() {
		data, _ := io.ReadAll(reader)
		output <- string(data)
	}()

	f()
	writer.Close()
	return <-output
}

func TestParseUploadResult(t *testing.T) {
	tests := []struct {
		name          string
		raw           string
		expected      distroUploadResult
		expectedError string
	}{
		{
			name:     "Test legacy success",
			raw:      `1`,
			expected: distroUploadResult{},
		},
		{
			name:          "Test legacy failure",
			raw:           `0`,
			expectedError: "unable to upload distro, server returned an error",
		},
		{
			name: "Test structured result",
			raw:  `{"path":"/srv/www/distributions/test","name":"test","registered":"test-label"}`,
			expected: distroUploadResult{
				Path: "/srv/www/distributions/test", Name: "test", Registered: "test-label", detection: true,
			},
		},
		{
			name:          "Test invalid result",
			raw:           `"garbage"`,
			expectedError: "unable to parse the distro upload result",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseUploadResult(json.RawMessage(tt.raw))
			if tt.expectedError != "" {
				testutils.AssertTrue(t, "parseUploadResult should fail", err != nil)
				testutils.AssertStringContains(t, "Unexpected error message", err.Error(), tt.expectedError)
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			testutils.AssertEquals(t, "Unexpected result", tt.expected, *result)
		})
	}
}

func TestDefaultDistroName(t *testing.T) {
	tests := []struct {
		filename string
		name     string
		expected string
	}{
		{filename: "test.iso", expected: "test"},
		{filename: "TEST.ISO", expected: "TEST"},
		{filename: "test", expected: "test"},
		{filename: "rhel-server-7.0-x86_64-dvd.iso", name: "rhel-server-7.0", expected: "rhel-server-7.0"},
	}

	for _, tt := range tests {
		t.Run(tt.filename+"/"+tt.name, func(t *testing.T) {
			testutils.AssertEquals(t, "Unexpected name", tt.expected, defaultDistroName(tt.filename, tt.name))
		})
	}
}

func TestValidateDistroName(t *testing.T) {
	tests := []struct {
		name          string
		expectedError string
	}{
		{name: "rhel-server-7.0"},
		{name: "a/b", expectedError: "must not contain path separators"},
		{name: "a\\b", expectedError: "must not contain path separators"},
		{name: ".hidden", expectedError: "must not start with a dot"},
		{name: " ", expectedError: "must not be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDistroName(tt.name)
			if tt.expectedError == "" {
				testutils.AssertTrue(t, "Unexpected error", err == nil)
				return
			}
			testutils.AssertTrue(t, "validateDistroName should fail", err != nil)
			testutils.AssertStringContains(t, "Unexpected error message", err.Error(), tt.expectedError)
		})
	}
}
