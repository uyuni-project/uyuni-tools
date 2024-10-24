// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// AddAPIFlags is a helper to include api details for the provided command tree.
func AddAPIFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("api-server", "", L("FQDN of the server to connect to"))
	cmd.PersistentFlags().String("api-user", "", L("API user username"))
	cmd.PersistentFlags().String("api-password", "", L("Password for the API user"))
	cmd.PersistentFlags().String("api-cacert", "", L("Path to a cert file of the CA"))
	cmd.PersistentFlags().Bool("api-insecure", false, L("If set, server certificate will not be checked for validity"))
}

func logTraceHeader(v *http.Header) {
	// Return early when not in trace loglevel
	if log.Logger.GetLevel() != zerolog.TraceLevel {
		return
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return
	}
	log.Trace().Msg(string(b))
}

func (c *APIClient) sendRequest(req *http.Request) (*http.Response, error) {
	log.Debug().Msgf("Sending %s request %s", req.Method, req.URL)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	if c.AuthCookie != nil {
		req.AddCookie(c.AuthCookie)
	}

	logTraceHeader(&req.Header)

	res, err := c.Client.Do(req)
	if err != nil {
		log.Trace().Err(err).Msgf("Request failed")
		return nil, err
	}

	logTraceHeader(&res.Header)

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		if res.StatusCode == 401 {
			return nil, errors.New(L("401: unauthorized"))
		}
		var errResponse map[string]string
		if res.Body != nil {
			body, err := io.ReadAll(res.Body)
			if err == nil {
				if err = json.Unmarshal(body, &errResponse); err == nil {
					error_message := fmt.Sprintf("%d: '%s'", res.StatusCode, errResponse["message"])
					return nil, errors.New(error_message)
				} else {
					error_message := fmt.Sprintf("%d: '%s'", res.StatusCode, string(body))
					return nil, errors.New(error_message)
				}
			}
		}
		return nil, fmt.Errorf(L("unknown error: %d"), res.StatusCode)
	}
	log.Debug().Msgf("Received response with code %d", res.StatusCode)

	return res, nil
}

// Init returns a HTTPClient object for further API use.
//
// Provided connectionDetails must have Server specified with FQDN to the
// target host.
//
// Optionaly connectionDetails can have user name and password set and Init
// will try to login to the host.
// caCert can be set to use custom CA certificate to validate target host.
func Init(conn *ConnectionDetails) (*APIClient, error) {
	// Load stored credentials as it also loads up server URL and CApath
	getStoredConnectionDetails(conn)

	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Warn().Msg(err.Error())
	}
	if conn.CApath != "" {
		caCert, err := os.ReadFile(conn.CApath)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	if conn.Server == "" {
		return nil, errors.New(L("server URL is not provided"))
	}
	client := &APIClient{
		Details: conn,
		BaseURL: fmt.Sprintf("https://%s%s", conn.Server, root_path_apiv1),
		Client: &http.Client{
			Timeout: time.Minute,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:            caCertPool,
					InsecureSkipVerify: conn.Insecure,
				},
			},
		},
	}
	if conn.Cookie != "" {
		client.AuthCookie = &http.Cookie{
			Name:  "pxt-session-cookie",
			Value: conn.Cookie,
		}
	}

	return client, err
}

// Login to the server using stored or provided credentials.
func (c *APIClient) Login() error {
	if c.Details.InSession {
		if err := c.sessionValidity(); err == nil {
			// Session is valid
			return nil
		}
		log.Warn().Msg(L("Cached session is expired."))
		if err := RemoveLoginCreds(); err != nil {
			log.Warn().Err(err).Msg(L("Failed to remove stored credentials!"))
		}
	}
	if err := getLoginCredentials(c.Details); err != nil {
		return err
	}
	return c.login()
}

func (c *APIClient) login() error {
	conn := c.Details
	url := fmt.Sprintf("%s/%s", c.BaseURL, "auth/login")
	data := map[string]string{
		"login":    conn.User,
		"password": conn.Password,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg(L("Unable to create login data"))
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	res, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}
	if !response["success"].(bool) {
		return fmt.Errorf(response["message"].(string))
	}

	cookies := res.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "pxt-session-cookie" && cookie.MaxAge > 0 {
			c.AuthCookie = cookie
			break
		}
	}

	if c.AuthCookie == nil {
		return errors.New(L("auth cookie not found in login response"))
	}

	return nil
}

func (c *APIClient) sessionValidity() error {
	// This is how spacecmd does it
	_, err := c.Get("user/listAssignableRoles")
	return err
}

// Logout from the server and remove localy stored session key.
func (c *APIClient) Logout() error {
	if _, err := c.Post("auth/logout", nil); err != nil {
		return utils.Errorf(err, L("failed to logout from the server"))
	}
	if err := RemoveLoginCreds(); err != nil {
		return err
	}
	return nil
}

// ValidateCreds checks if the login credentials are valid.
func (c *APIClient) ValidateCreds() bool {
	err := c.Login()
	return err == nil
}

// Post issues a POST HTTP request to the API target
//
// `path` specifies an API endpoint
// `data` contains a map of values to add to the POST query. `data` are serialized to the JSON
//
// returns a raw HTTP Response.
func (c *APIClient) Post(path string, data map[string]interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, path)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg(L("Unable to convert data to JSON"))
		return nil, err
	}

	log.Trace().Msgf("payload: %s", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get issues GET HTTP request to the API target
//
// `path` specifies API endpoint together with query options
//
// returns a raw HTTP Response.
func (c *APIClient) Get(path string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Post issues a POST HTTP request to the API target using the client and decodes the response.
//
// `path` specifies an API endpoint
// `data` contains a map of values to add to the POST query. `data` are serialized to the JSON
//
// returns a deserialized JSON data to the map.
func Post[T interface{}](client *APIClient, path string, data map[string]interface{}) (*ApiResponse[T], error) {
	res, err := client.Post(path, data)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var response ApiResponse[T]
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Trace().Msgf("response: %s", string(body))

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Get issues an HTTP GET request to the API using the client and decodes the response.
//
// `path` specifies API endpoint together with query options
//
// returns an ApiResponse with the decoded result.
func Get[T interface{}](client *APIClient, path string) (*ApiResponse[T], error) {
	res, err := client.Get(path)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var response ApiResponse[T]
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
