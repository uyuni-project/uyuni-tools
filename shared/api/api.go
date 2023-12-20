// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const ROOT_PATH_APIv1 = "/rhn/manager/api"

// HTTP Client is an API entrypoint
type HTTPClient struct {

	// URL to the API endpoint of the target host
	BaseURL string

	// net/http client
	Client *http.Client

	// Authentication cookie storage
	AuthCookie *http.Cookie
}

// Connection details for initial API connection
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

// AddAPIFlags is a helper to include api details for the provided command tree.
//
// If the api support is only optional for the command, set optional parameter to true
func AddAPIFlags(cmd *cobra.Command, optional bool) {
	cmd.PersistentFlags().String("api-server", "", "FQDN of the server to connect to")
	cmd.PersistentFlags().String("api-user", "", "API user username")
	cmd.PersistentFlags().String("api-password", "", "Password for the API user")
	cmd.PersistentFlags().String("api-cacert", "", "Path to a cert file of the CA")
	cmd.PersistentFlags().Bool("api-insecure", false, "If set, server certificate will not be checked for validity")

	if !optional {
		cmd.MarkFlagRequired("api-server")
		cmd.MarkFlagRequired("api-username")
		cmd.MarkFlagRequired("api-password")
	}
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
		return nil, err
	}

	log.Trace().Msg(prettyPrint(res.Header))
	log.Trace().Msg(prettyPrint(res.Body))

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errResponse map[string]string
		if err = json.NewDecoder(res.Body).Decode(&errResponse); err == nil {
			return nil, fmt.Errorf(errResponse["message"])
		}
		return nil, fmt.Errorf("unknown error: %d", res.StatusCode)
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
// caCert can be set to use custom CA certificate to validate target host
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
		BaseURL: fmt.Sprintf("https://%s%s", conn.Server, ROOT_PATH_APIv1),
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
		log.Error().Err(err).Msg("Unable to create login data")
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
		return fmt.Errorf("auth cookie not found in login response")
	}

	return nil
}

// Post issues a POST HTTP request to the API target
//
// `path` specifies an API endpoint
// `data` contains a map of values to add to the POST query. `data` are serialized to the JSON
//
// returns a deserialized JSON data to the map
func (c *HTTPClient) Post(path string, data map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, path)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("Unable to JSONify data")
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

// Get issues GET HTTP request to the API target
//
// `path` specifies API endpoint together with query options
//
// returns a deserialized JSON data to the map
func (c *HTTPClient) Get(path string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}
