// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/mocks"
	"github.com/uyuni-project/uyuni-tools/shared/api/proxy"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

// global access to the testing object.
var globalT *testing.T

// ProxyConfigGenerateRequestBodyData is the data structure for the body of the ContainerConfigGenerate API request.
type ProxyConfigGenerateRequestBodyData struct {
	ProxyName  string
	ProxyPort  int
	Server     string
	MaxCache   int
	Email      string
	CaCrt      string
	CaKey      string
	CaPassword string
	Cnames     []string
	Country    string
	State      string
	City       string
	Org        string
	OrgUnit    string
	SSLEmail   string
}

// ProxyConfigRequestBodyData is the data structure for the request body of the ContainerConfig API request.
type ProxyConfigRequestBodyData struct {
	ProxyName       string
	ProxyPort       int
	Server          string
	MaxCache        int
	Email           string
	IntermediateCAs []string
	ProxyCrt        string
	ProxyKey        string
	RootCA          string
}

// common connection details (for generating the client).
const user = "testUser"
const password = "testPwd"
const server = "testServer"

var connectionDetails = &api.ConnectionDetails{User: user, Password: password, Server: server}

// common expected values for both ContainerConfig and ContainerConfigGenerate calls.
const expectedProxyName = "testProxy"
const expectedProxyPort = 8080
const expectedServer = "testServer"
const expectedMaxCache = 100
const expectedEmail = "test@email.com"

// expected values for ContainerConfig.
const expectedCaCrt = "caCrt contents"
const expectedCaKey = "caKey contents"
const expectedCaPassword = "caPwd"
const expectedCountry = "testCountry"
const expectedState = "exampleState"
const expectedCity = "exampleCity"
const expectedOrg = "exampleOrg"
const expectedOrgUnit = "exampleOrgUnit"
const expectedSSLEmail = "sslEmail@example.com"

var expectedCnames = []string{"altNameA.example.com", "altNameB.example.com"}

// expected values for ContainerConfigGenerate.
const expectedRootCA = "rootCA contents"
const expectedProxyCrt = "proxyCrt contents"
const expectedProxyKey = "proxyKey contents"

var expectedIntermediateCAs = []string{"intermediateCA1", "intermediateCA2"}

var proxyConfigRequest = proxy.ProxyConfigRequest{
	ProxyName:       expectedProxyName,
	ProxyPort:       expectedProxyPort,
	Server:          expectedServer,
	MaxCache:        expectedMaxCache,
	Email:           expectedEmail,
	RootCA:          expectedRootCA,
	ProxyCrt:        expectedProxyCrt,
	ProxyKey:        expectedProxyKey,
	IntermediateCAs: expectedIntermediateCAs,
}

var proxyConfigGenerateRequest = proxy.ProxyConfigGenerateRequest{
	ProxyName:  expectedProxyName,
	ProxyPort:  expectedProxyPort,
	Server:     expectedServer,
	MaxCache:   expectedMaxCache,
	Email:      expectedEmail,
	CaCrt:      expectedCaCrt,
	CaKey:      expectedCaKey,
	CaPassword: expectedCaPassword,
	Cnames:     expectedCnames,
	Country:    expectedCountry,
	State:      expectedState,
	City:       expectedCity,
	Org:        expectedOrg,
	OrgUnit:    expectedOrgUnit,
	SSLEmail:   expectedSSLEmail,
}

// Tests ContainerConfig when the post request fails.
func TestFailContainerConfigWhenPostRequestFails(t *testing.T) {
	//
	expectedErrorMessage := "failed to create proxy configuration file"

	// Mock client
	client, err := api.Init(connectionDetails)
	if err != nil {
		t.FailNow()
	}
	client.Client = &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutils.GetResponse(404, `{}`)
		},
	}

	// Execute
	result, err := proxy.ContainerConfig(client, proxyConfigRequest)

	// Assertions
	testutils.AssertTrue(t, "Unexpected successful ContainerConfigGenerate call", err != nil)
	testutils.AssertTrue(t, "ContainerConfigGenerate error message", strings.Contains(err.Error(), expectedErrorMessage))
	testutils.AssertTrue(t, "Result data should be nil", result == nil)
}

// Tests ContainerConfig when the post request is successful but the response is unsuccessful.
func TestFailContainerConfigWhenPostIsUnsuccessful(t *testing.T) {
	//
	expectedErrorMessage := "some error message"

	// Mock client
	client, err := api.Init(connectionDetails)
	if err != nil {
		t.FailNow()
	}
	client.Client = &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutils.GetResponse(200, `{"success": false, "message": "some error message"}`)
		},
	}

	// Execute
	result, err := proxy.ContainerConfig(client, proxyConfigRequest)

	// Assertions
	testutils.AssertTrue(t, "Unexpected successful ContainerConfigGenerate call", err != nil)
	testutils.AssertTrue(t, "ContainerConfigGenerate error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	testutils.AssertTrue(t, "Result data should be nil", result == nil)
}

// Tests ContainerConfig when all parameters are provided.
func TestSuccessfulContainerConfigWhenAllParametersAreProvided(t *testing.T) {
	//
	expectedResponseData := []int8{1, 2, 3, 4, 5}

	// Mock client
	client, err := api.Init(connectionDetails)
	if err != nil {
		t.FailNow()
	}
	client.Client = &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// asserts request body contents
			var data ProxyConfigRequestBodyData
			if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
				return nil, err
			}

			testutils.AssertEquals(globalT, "ProxyName doesn't match", expectedProxyName, data.ProxyName)
			testutils.AssertEquals(globalT, "ProxyPort doesn't match", expectedProxyPort, data.ProxyPort)
			testutils.AssertEquals(globalT, "Server doesn't match", expectedServer, data.Server)
			testutils.AssertEquals(globalT, "MaxCache doesn't match", expectedMaxCache, data.MaxCache)
			testutils.AssertEquals(globalT, "Email doesn't match", expectedEmail, data.Email)
			testutils.AssertEquals(globalT, "RootCA doesn't match", expectedRootCA, data.RootCA)
			testutils.AssertEquals(globalT, "ProxyCrt doesn't match", expectedProxyCrt, data.ProxyCrt)
			testutils.AssertEquals(globalT, "ProxyKey doesn't match", expectedProxyKey, data.ProxyKey)
			testutils.AssertEquals(globalT, "intermediateCas don't match",
				fmt.Sprintf("%v", expectedIntermediateCAs),
				fmt.Sprintf("%v", data.IntermediateCAs))

			// mock response
			return testutils.GetResponse(200, `{"success": true, "result": [1, 2, 3, 4, 5]}`)
		},
	}

	// Execute
	result, err := proxy.ContainerConfig(client, proxyConfigRequest)

	// Assertions
	testutils.AssertTrue(t, "Unexpected error executing ContainerConfigGenerate", err == nil)
	testutils.AssertTrue(t, "Result should not be empty", result != nil)
	testutils.AssertEquals(
		t, "Result configuration binary doesn't match",
		fmt.Sprintf("%v", expectedResponseData), fmt.Sprintf("%v", *result),
	)
}

// Tests ContainerConfigGenerate when the post request fails.
func TestFailContainerConfigGenerateWhenPostRequestFails(t *testing.T) {
	//
	expectedErrorMessage := "failed to create proxy configuration file"

	// Mock client
	client, err := api.Init(connectionDetails)
	if err != nil {
		t.FailNow()
	}
	client.Client = &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutils.GetResponse(404, `{}`)
		},
	}

	// Execute
	result, err := proxy.ContainerConfigGenerate(client, proxyConfigGenerateRequest)

	// Assertions
	testutils.AssertTrue(t, "Unexpected successful ContainerConfigGenerate call", err != nil)
	testutils.AssertTrue(t, "ContainerConfigGenerate error message", strings.Contains(err.Error(), expectedErrorMessage))
	testutils.AssertTrue(t, "Result data should be nil", result == nil)
}

// Tests ContainerConfigGenerate when the post request is successful but the response is unsuccessful.
func TestFailContainerConfigGenerateWhenPostIsUnsuccessful(t *testing.T) {
	//
	expectedErrorMessage := "some error message"

	// Mock client
	client, err := api.Init(connectionDetails)
	if err != nil {
		t.FailNow()
	}
	client.Client = &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutils.GetResponse(200, `{"success": false, "message": "some error message"}`)
		},
	}

	// Execute
	result, err := proxy.ContainerConfigGenerate(client, proxyConfigGenerateRequest)

	// Assertions
	testutils.AssertTrue(t, "Unexpected successful ContainerConfigGenerate call", err != nil)
	testutils.AssertTrue(t, "ContainerConfigGenerate error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	testutils.AssertTrue(t, "Result data should be nil", result == nil)
}

// Tests ContainerConfig when all parameters are provided.
func TestSuccessfulContainerConfigGenerateWhenAllParametersAreProvided(t *testing.T) {
	//
	expectedResponseData := []int8{1, 2, 3, 4, 5}

	// Mock client
	client, err := api.Init(connectionDetails)
	if err != nil {
		t.FailNow()
	}
	client.Client = &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// asserts request body contents
			var data ProxyConfigGenerateRequestBodyData
			if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
				return nil, err
			}

			testutils.AssertEquals(globalT, "ProxyName doesn't match", expectedProxyName, data.ProxyName)
			testutils.AssertEquals(globalT, "ProxyPort doesn't match", expectedProxyPort, data.ProxyPort)
			testutils.AssertEquals(globalT, "Server doesn't match", expectedServer, data.Server)
			testutils.AssertEquals(globalT, "MaxCache doesn't match", expectedMaxCache, data.MaxCache)
			testutils.AssertEquals(globalT, "Email doesn't match", expectedEmail, data.Email)

			testutils.AssertEquals(globalT, "CaCertificate doesn't match", expectedCaCrt, data.CaCrt)
			testutils.AssertEquals(globalT, "CaKey doesn't match", expectedCaKey, data.CaKey)
			testutils.AssertEquals(globalT, "CaPassword doesn't match", expectedCaPassword, data.CaPassword)
			testutils.AssertEquals(
				globalT, "Cnames don't match", fmt.Sprintf("%v", expectedCnames), fmt.Sprintf("%v", data.Cnames),
			)
			testutils.AssertEquals(globalT, "Country doesn't match", expectedCountry, data.Country)
			testutils.AssertEquals(globalT, "State doesn't match", expectedState, data.State)
			testutils.AssertEquals(globalT, "City doesn't match", expectedCity, data.City)
			testutils.AssertEquals(globalT, "Org doesn't match", expectedOrg, data.Org)
			testutils.AssertEquals(globalT, "OrgUnit doesn't match", expectedOrgUnit, data.OrgUnit)
			testutils.AssertEquals(globalT, "SSLEmail doesn't match", expectedSSLEmail, data.SSLEmail)

			// mock response
			return testutils.GetResponse(200, `{"success": true, "result": [1, 2, 3, 4, 5]}`)
		},
	}

	// Execute
	result, err := proxy.ContainerConfigGenerate(client, proxyConfigGenerateRequest)

	// Assertions
	testutils.AssertTrue(t, "Unexpected error executing ContainerConfigGenerate", err == nil)
	testutils.AssertTrue(t, "Result should not be empty", result != nil)
	testutils.AssertEquals(
		t, "Result configuration binary doesn't match",
		fmt.Sprintf("%v", expectedResponseData), fmt.Sprintf("%v", *result),
	)
}
