// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import "net/http"

const rootPathApiv1 = "/rhn/manager/api"
const apiCredentialsStore = ".uyuni-api.json"

// APIClient is the API entrypoint structure.
type APIClient struct {

	// URL to the API endpoint of the target host
	BaseURL string

	// net/http client
	Client HTTPClient

	// Authentication cookie storage
	AuthCookie *http.Cookie

	// Connection details
	Details *ConnectionDetails

	// Valid API namespaces advertised by the server
	validNamespaces map[string]struct{}
}

// NamespaceAccess describes an API namespace advertised by the server.
type NamespaceAccess struct {
	Namespace   string `json:"namespace"`
	Description string `json:"description"`
	AccessMode  string `json:"access_mode"`
}

// HTTPClient is a minimal HTTPClient interface primarily for unit testing.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ConnectionDetails holds the details for initial API connection.
type ConnectionDetails struct {

	// FQDN of the target host.
	Server string

	// User to login under.
	User string

	// Password for the user.
	Password string

	// Path to CA certificate file used for target host validation.
	// Provided certificate is used together with system certificates.
	CApath string `mapstructure:"cacert"`

	// Disable certificate validation, unsecure and not recommended.
	Insecure bool

	// Indicates if details we loaded from cache
	InSession bool

	// PXE cookie
	Cookie string
}

// APIResponse describes the HTTP response where T is the type of the result.
type APIResponse[T interface{}] struct {
	Result  T
	Success bool
	Message string
}

// Authentication storage.
type authStorage struct {
	Session string
	Server  string
	CApath  string
}
