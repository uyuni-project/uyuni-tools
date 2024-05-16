// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const root_path_apiv1 = "/rhn/manager/api"

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

// AddAPIFlags is a helper to include api details for the provided command tree.
//
// If the api support is only optional for the command, set optional parameter to true.
func AddAPIFlags(cmd *cobra.Command, optional bool) error {
	cmd.PersistentFlags().String("api-server", "", L("FQDN of the server to connect to"))
	cmd.PersistentFlags().String("api-user", "", L("API user username"))
	cmd.PersistentFlags().String("api-password", "", L("Password for the API user"))
	cmd.PersistentFlags().String("api-cacert", "", L("Path to a cert file of the CA"))
	cmd.PersistentFlags().Bool("api-insecure", false, L("If set, server certificate will not be checked for validity"))

	if !optional {
		if err := cmd.MarkPersistentFlagRequired("api-server"); err != nil {
			return err
		}
		if err := cmd.MarkPersistentFlagRequired("api-user"); err != nil {
			return err
		}
		if err := cmd.MarkPersistentFlagRequired("api-password"); err != nil {
			return err
		}
	}
	return nil
}

func prettyPrint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return fmt.Sprintln(string(b))
}

func (c *HTTPClient) sendRequest(req *http.Request) (*http.Response, error) {
	log.Debug().Msgf("Sending %s request %s", req.Method, req.URL)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	if c.AuthCookie != nil {
		req.AddCookie(c.AuthCookie)
	}

	log.Trace().Msg(prettyPrint(req.Header))
	log.Trace().Msg(prettyPrint(req.Body))

	res, err := c.Client.Do(req)
	if err != nil {
		log.Trace().Err(err).Msgf("Request failed")
		return nil, err
	}

	log.Trace().Msg(prettyPrint(res.Header))
	log.Trace().Msg(prettyPrint(res.Body))

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errResponse map[string]string
		if err = json.NewDecoder(res.Body).Decode(&errResponse); err == nil {
			return nil, fmt.Errorf(errResponse["message"])
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
func Init(conn *ConnectionDetails) (*HTTPClient, error) {
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Warn().Msg(err.Error())
	}
	if conn.CAcert != "" {
		caCert, err := os.ReadFile(conn.CAcert)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}
	client := &HTTPClient{
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

	if len(conn.User) > 0 {
		if len(conn.Password) == 0 {
			utils.AskPasswordIfMissing(&conn.Password, L("API server password"), 0, 0)
		}
		err = client.login(conn)
	}
	return client, err
}

func (c *HTTPClient) login(conn *ConnectionDetails) error {
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
		return fmt.Errorf(response["messages"].(string))
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

// Post issues a POST HTTP request to the API target
//
// `path` specifies an API endpoint
// `data` contains a map of values to add to the POST query. `data` are serialized to the JSON
//
// returns a raw HTTP Response.
func (c *HTTPClient) Post(path string, data map[string]interface{}) (*http.Response, error) {
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
func (c *HTTPClient) Get(path string) (*http.Response, error) {
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
func Post[T interface{}](client *HTTPClient, path string, data map[string]interface{}) (*ApiResponse[T], error) {
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
func Get[T interface{}](client *HTTPClient, path string) (*ApiResponse[T], error) {
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
