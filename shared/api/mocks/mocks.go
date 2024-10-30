// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package mocks

import "net/http"

// MockClient is a mocked api.HTTPClient.
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do fulfills the api.HTTPClient interface.
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
