// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

var mockSystemJSON = `{
	"success": true,
	"result": [
		{
			"id": 1001,
			"name": "test-system.uy",
			"last_checkin": "2026-03-19 14:00:00",
			"created": "2026-03-19 13:00:00"
		}
	]
}`

func TestRunSystem_JSON(t *testing.T) {
	// 1. Create a mocked TLS server that simulates the Uyuni API
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock the login authentication endpoint
		if strings.HasSuffix(r.URL.Path, "auth/login") {
			w.Header().Set("Set-Cookie", "pxt-session-cookie=mock_cookie; Max-Age=3600")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"success": true, "result": "mock_cookie"}`))
			return
		}

		// Mock the actual system endpoint
		if strings.HasSuffix(r.URL.Path, "system/listSystems") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(mockSystemJSON))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	// 2. Configure the mocked API Server details
	mockHost := strings.TrimPrefix(ts.URL, "https://")

	globalFlags := &types.GlobalFlags{}
	flags := &getFlags{
		ConnectionDetails: api.ConnectionDetails{
			Server:   mockHost,
			User:     "admin",
			Password: "mockpassword",
			Insecure: true,
		},
		OutputFormat: "json",
	}

	// 3. Execute the function natively
	cmd := newSystemCommand(globalFlags, flags)
	err := cmd.RunE(cmd, []string{})

	// 4. Assert no error occurred during authentication, fetching, and formatting
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
