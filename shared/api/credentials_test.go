// SPDX-FileCopyrightText: 2024 SUSE LLC
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

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/api/mocks"
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
	if err != nil {
		t.FailNow()
	}

	client.Client = &mocks.MockClient{
		DoFunc: loginTestDo,
	}

	if err := client.Login(); err != nil {
		t.FailNow()
	}
	err = StoreLoginCreds(client)
	if err != nil {
		t.Fail()
	}

	connection2 := ConnectionDetails{}
	if err := loadLoginCreds(&connection2); err != nil {
		t.Fail()
	}
	if connection2.Server != server {
		log.Error().Msg("server does not match")
		t.Fail()
	}
	if connection2.Cookie != cookie {
		log.Error().Msg("cookie does not match")
		t.Fail()
	}
}

// Test credentials are cleaned-up after logout.
func TestCredentialsCleanup(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := storeTestCredentials()
	if err != nil {
		log.Error().Err(err).Msg("failed to store creds")
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
	if err := RemoveLoginCreds(); err != nil {
		t.Fail()
	}

	connection2 := ConnectionDetails{
		Server: server,
	}
	err = loadLoginCreds(&connection2)
	if err == nil {
		t.Fail()
	}
	if connection2.User != "" {
		t.Fail()
	}
}

// Write malformed credentials file to check autocleanup of wrong credentials.
func TestAutocleanup(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := os.WriteFile(getAPICredsFile(), []byte(""), 0600)
	if err != nil {
		t.Fail()
	}
	connection := ConnectionDetails{
		Server: server,
	}
	getStoredConnectionDetails(&connection)
	if connection.InSession {
		t.Fail()
	}
	_, err = os.Stat(getAPICredsFile())
	if err == nil {
		t.Fail()
	}
}

// Test login using cached credentials.
func TestCredentialValidation(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := storeTestCredentials()
	if err != nil {
		log.Error().Err(err).Msg("failed to store creds")
		t.Fail()
	}

	connection := ConnectionDetails{}
	client, err := Init(&connection)
	if err != nil {
		log.Error().Err(err).Msg("failed to init connection")
		t.Fail()
	}

	if !client.Details.InSession {
		log.Error().Msg("Credentials are not marked as cached")
		t.Fail()
	}

	client.Client = &mocks.MockClient{
		DoFunc: userListRolesDo,
	}

	if err := client.Login(); err != nil {
		log.Trace().Err(err).Msg("failed")
		t.Fail()
	}
}

func TestWrongCredentials(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := storeWrongTestCredentials()
	if err != nil {
		log.Error().Err(err).Msg("failed to store creds")
		t.Fail()
	}

	connection := ConnectionDetails{}
	client, err := Init(&connection)
	if err != nil {
		log.Error().Err(err).Msg("failed to init connection")
		t.Fail()
	}

	if !client.Details.InSession {
		log.Error().Msg("Credentials are not marked as cached")
		t.Fail()
	}

	client.Client = &mocks.MockClient{
		DoFunc: userListRolesDo,
	}

	err = client.Login()
	if err == nil {
		log.Error().Err(err).Msg("login was successful even when should not have been")
		t.Fail()
	}

	// Test that wrong login will remove auth file
	_, err = os.Stat(getAPICredsFile())
	if err == nil {
		t.Fail()
	}
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
	headers.Add("Set-Cookie", fmt.Sprintf("pxt-session-cookie=%s; Max-Age=3600; Path=/; Secure; HttpOnly;HttpOnly;Secure", cookie))
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
	headers.Add("Set-Cookie", fmt.Sprintf("pxt-session-cookie=%s; Max-Age=3600; Path=/; Secure; HttpOnly;HttpOnly;Secure", cookie))
	return &http.Response{
		StatusCode: 200,
		Header:     headers,
		Body:       r,
	}, nil
}
