// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/api/mocks"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

const user = "mytestuser"
const password = "mytestpassword"
const server = "mytestserver"
const cookie = "mytestpxtcookie"

// Test happy path for credentials store.
func TestCredentialsStore(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	connection := ConnectionDetails{
		User:     user,
		Password: password,
		Server:   server,
	}
	client, err := Init(&connection)
	testutils.AssertNoError(t, "api.Init failed", err)

	client.Client = &mocks.MockClient{
		DoFunc: loginTestDo,
	}

	testutils.AssertNoError(t, "client.Login failed", client.Login())
	testutils.AssertNoError(t, "StoreLoginCreds failed", StoreLoginCreds(client))

	connection2 := ConnectionDetails{}
	testutils.AssertNoError(t, "loadLoginCreds failed", loadLoginCreds(&connection2))
	testutils.AssertEquals(t, "server does not match", server, connection2.Server)
	testutils.AssertEquals(t, "cookie does not match", cookie, connection2.Cookie)
}

// Test credentials are cleaned-up after logout.
func TestCredentialsCleanup(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := storeTestCredentials()
	testutils.AssertNoError(t, "failed to store creds", err)

	testutils.AssertNoError(t, "RemoveLoginCreds failed", RemoveLoginCreds())

	connection2 := ConnectionDetails{
		Server: server,
	}
	err = loadLoginCreds(&connection2)
	testutils.AssertError(t, "unable to read credentials file", err)
	testutils.AssertEquals(t, "expected empty user", "", connection2.User)
}

// Write malformed credentials file to check autocleanup of wrong credentials.
func TestAutocleanup(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := os.WriteFile(getAPICredsFile(), []byte(""), 0600)
	testutils.AssertNoError(t, "WriteFile failed", err)
	connection := ConnectionDetails{
		Server: server,
	}
	getStoredConnectionDetails(&connection)
	testutils.AssertTrue(t, "connection InSession should not be set", !connection.InSession)
	_, err = os.Stat(getAPICredsFile())
	testutils.AssertError(t, "no such file or directory", err)
}

// Test login using cached credentials.
func TestCredentialValidation(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := storeTestCredentials()
	testutils.AssertNoError(t, "failed to store creds", err)

	connection := ConnectionDetails{}
	client, err := Init(&connection)
	testutils.AssertNoError(t, "failed to init connection", err)
	testutils.AssertTrue(t, "Credentials are not marked as cached", client.Details.InSession)

	client.Client = &mocks.MockClient{
		DoFunc: userListRolesDo,
	}

	testutils.AssertNoError(t, "client.Login failed", client.Login())
}

func TestWrongCredentials(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := storeWrongTestCredentials()
	testutils.AssertNoError(t, "failed to store creds", err)

	connection := ConnectionDetails{}
	client, err := Init(&connection)
	testutils.AssertNoError(t, "failed to init connection", err)
	testutils.AssertTrue(t, "Credentials are not marked as cached", client.Details.InSession)

	client.Client = &mocks.MockClient{
		DoFunc: userListRolesDo,
	}

	err = client.Login()
	testutils.AssertError(t, "No credentials provided", err)

	// Test that wrong login will remove auth file
	_, err = os.Stat(getAPICredsFile())
	testutils.AssertError(t, "no such file or directory", err)
}

// helper storing valid credentials.
func storeTestCredentials() error {
	client := APIClient{
		Details: &ConnectionDetails{
			User:   user,
			Server: server,
		},
		AuthCookie: &http.Cookie{
			Name:  "pxt-session-cookie",
			Value: cookie,
		},
	}
	return StoreLoginCreds(&client)
}

// helper storing invalid credentials.
func storeWrongTestCredentials() error {
	client := APIClient{
		Details: &ConnectionDetails{
			User:   user,
			Server: server,
		},
		AuthCookie: &http.Cookie{
			Name:  "pxt-session-cookie",
			Value: "wrongcookie",
		},
	}
	return StoreLoginCreds(&client)
}

// helper login request response.
func loginTestDo(req *http.Request) (*http.Response, error) {
	if req.URL.Path != "/rhn/manager/api/auth/login" {
		return &http.Response{
			StatusCode: 404,
		}, nil
	}
	data := map[string]string{}
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		return nil, err
	}
	if data["login"] != user || data["password"] != password {
		return &http.Response{
			StatusCode: 403,
		}, nil
	}
	json := `{"success": true}`
	r := io.NopCloser(bytes.NewReader([]byte(json)))
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add(
		"Set-Cookie",
		fmt.Sprintf("pxt-session-cookie=%s; Max-Age=3600; Path=/; Secure; HttpOnly;HttpOnly;Secure", cookie),
	)
	return &http.Response{
		StatusCode: 200,
		Header:     headers,
		Body:       r,
	}, nil
}

func userListRolesDo(req *http.Request) (*http.Response, error) {
	if req.URL.Path != "/rhn/manager/api/user/listAssignableRoles" {
		return &http.Response{
			StatusCode: 404,
		}, nil
	}
	if pxt, err := req.Cookie("pxt-session-cookie"); err != nil || pxt.Value != cookie {
		return &http.Response{
			StatusCode: 403,
		}, nil
	}
	json := `{"success": true}`
	r := io.NopCloser(bytes.NewReader([]byte(json)))
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add(
		"Set-Cookie",
		fmt.Sprintf("pxt-session-cookie=%s; Max-Age=3600; Path=/; Secure; HttpOnly;HttpOnly;Secure", cookie),
	)
	return &http.Response{
		StatusCode: 200,
		Header:     headers,
		Body:       r,
	}, nil
}
