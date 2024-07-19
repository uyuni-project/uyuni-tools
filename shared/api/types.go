// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import "net/http"

const root_path_apiv1 = "/rhn/manager/api"
const api_credentials_store = ".uyuni-api.json"

// HTTP Client is an API entrypoint.
type HTTPClient struct {

	// URL to the API endpoint of the target host
	BaseURL string

	// net/http client
	Client *http.Client

	// Authentication cookie storage
	AuthCookie *http.Cookie
}

// Connection details for initial API connection.
type ConnectionDetails struct {

	// FQDN of the target host.
	Server string

	// User to login under.
	User string

	// Password for the user.
	Password string

	// CA certificate used for target host validation.
	// Provided certificate is used together with system certificates.
	CAcert string

	// Disable certificate validation, unsecure and not recommended.
	Insecure bool
}

// API response where T is the type of the result.
type ApiResponse[T interface{}] struct {
	Result  T
	Success bool
	Message string
}

// Authentication storage.
type authStorage struct {
	User     []byte
	Password []byte
	Server   string
}
