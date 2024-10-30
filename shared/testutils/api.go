// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const defaultCookie = "testCookie"

// GetResponse is a helper function to generate a response with the given status and json body.
func GetResponse(status int, json string) (*http.Response, error) {
	return GetResponseWithCookie(defaultCookie, status, json)
}

// GetResponseWithCookie is a helper function to generate a response with the given status and json body.
func GetResponseWithCookie(cookie string, status int, json string) (*http.Response, error) {
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add(
		"Set-Cookie",
		fmt.Sprintf("pxt-session-cookie=%s; Max-Age=3600; Path=/; Secure; HttpOnly;HttpOnly;Secure", cookie),
	)
	return &http.Response{
		StatusCode: status,
		Header:     headers,
		Body:       body,
	}, nil
}

// SuccessfulLoginTestDo is a helper function to mock a successful login response.
func SuccessfulLoginTestDo(req *http.Request) (*http.Response, error) {
	if req.URL.Path != "/rhn/manager/api/auth/login" {
		return &http.Response{
			StatusCode: 404,
		}, nil
	}

	return GetResponse(200, `{"success": true}`)
}

// FailedLoginTestDo is a helper function to mock a failed login response due to incorrect credentials.
func FailedLoginTestDo(req *http.Request) (*http.Response, error) {
	if req.URL.Path != "/rhn/manager/api/auth/login" {
		return &http.Response{
			StatusCode: 404,
		}, nil
	}

	return GetResponse(200, `{"success":false,"message":"Either the password or username is incorrect."}`)
}
