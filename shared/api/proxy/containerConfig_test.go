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
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
)

// global access to the testing object.
var globalT *testing.T

// ProxyConfigGenerateRequestBodyData is the data structure for the request body of the ContainerConfigGenerate API request.
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
	SslEmail   string
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
const expectedSslEmail = "sslEmail@example.com"

var expectedCnames = []string{"altNameA.example.com", "altNameB.example.com"}

// expected values for ContainerConfigGenerate.
const expectedRootCA = "rootCA contents"
const expectedProxyCrt = "proxyCrt contents"
const expectedProxyKey = "proxyKey contents"

var expectedIntermediateCAs = []string{"intermediateCA1", "intermediateCA2"}

// Tests ContainerConfig when the post request fails.
func TestFailContainerConfigWhenPostRequestFails(t *testing.T) {
	//
	expectedErrorMessage := "failed to create proxy configuration file with generated certificates"

	// Mock client
	client, err := api.Init(connectionDetails)
	if err != nil {
		t.FailNow()
	}
	client.Client = &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return test_utils.GetResponse(404, `{}`)
		},
	}
	// Execute
	result, err := proxy.ContainerConfig(client, expectedProxyName, expectedProxyPort,
		expectedServer, expectedMaxCache, expectedEmail,
		expectedRootCA, expectedProxyCrt, expectedProxyKey, expectedIntermediateCAs)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected successful ContainerConfigGenerate call", err != nil)
	test_utils.AssertTrue(t, "ContainerConfigGenerate error message", strings.Contains(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "Result data should be nil", result == nil)
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
			return test_utils.GetResponse(200, `{"success": false, "message": "some error message"}`)
		},
	}
	// Execute
	result, err := proxy.ContainerConfig(client, expectedProxyName, expectedProxyPort,
		expectedServer, expectedMaxCache, expectedEmail,
		expectedRootCA, expectedProxyCrt, expectedProxyKey, expectedIntermediateCAs)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected successful ContainerConfigGenerate call", err != nil)
	test_utils.AssertTrue(t, "ContainerConfigGenerate error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "Result data should be nil", result == nil)
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

			test_utils.AssertEquals(globalT, "ProxyName doesn't match", expectedProxyName, data.ProxyName)
			test_utils.AssertEquals(globalT, "ProxyPort doesn't match", expectedProxyPort, data.ProxyPort)
			test_utils.AssertEquals(globalT, "Server doesn't match", expectedServer, data.Server)
			test_utils.AssertEquals(globalT, "MaxCache doesn't match", expectedMaxCache, data.MaxCache)
			test_utils.AssertEquals(globalT, "Email doesn't match", expectedEmail, data.Email)
			test_utils.AssertEquals(globalT, "RootCA doesn't match", expectedRootCA, data.RootCA)
			test_utils.AssertEquals(globalT, "ProxyCrt doesn't match", expectedProxyCrt, data.ProxyCrt)
			test_utils.AssertEquals(globalT, "ProxyKey doesn't match", expectedProxyKey, data.ProxyKey)
			test_utils.AssertEquals(globalT, "intermediateCas don't match",
				fmt.Sprintf("%v", expectedIntermediateCAs),
				fmt.Sprintf("%v", data.IntermediateCAs))

			// mock response
			return test_utils.GetResponse(200, `{"success": true, "result": [1, 2, 3, 4, 5]}`)
		},
	}
	// Execute
	result, err := proxy.ContainerConfig(client, expectedProxyName, expectedProxyPort,
		expectedServer, expectedMaxCache, expectedEmail,
		expectedRootCA, expectedProxyCrt, expectedProxyKey, expectedIntermediateCAs)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected error executing ContainerConfigGenerate", err == nil)
	test_utils.AssertTrue(t, "Result should not be empty", result != nil)
	test_utils.AssertEquals(t, "Result configuration binary doesn't match", fmt.Sprintf("%v", expectedResponseData), fmt.Sprintf("%v", *result))
}

// Tests ContainerConfigGenerate when the post request fails.
func TestFailContainerConfigGenerateWhenPostRequestFails(t *testing.T) {
	//
	expectedErrorMessage := "failed to create proxy configuration file with generated certificates"

	// Mock client
	client, err := api.Init(connectionDetails)
	if err != nil {
		t.FailNow()
	}
	client.Client = &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return test_utils.GetResponse(404, `{}`)
		},
	}
	// Execute
	result, err := proxy.ContainerConfigGenerate(client, expectedProxyName, expectedProxyPort,
		expectedServer, expectedMaxCache, expectedEmail,
		expectedCaCrt, expectedCaKey, expectedCaPassword, expectedCnames, expectedCountry,
		expectedState, expectedCity, expectedOrg, expectedOrgUnit, expectedSslEmail)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected successful ContainerConfigGenerate call", err != nil)
	test_utils.AssertTrue(t, "ContainerConfigGenerate error message", strings.Contains(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "Result data should be nil", result == nil)
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
			return test_utils.GetResponse(200, `{"success": false, "message": "some error message"}`)
		},
	}
	// Execute
	result, err := proxy.ContainerConfigGenerate(client, expectedProxyName, expectedProxyPort,
		expectedServer, expectedMaxCache, expectedEmail,
		expectedCaCrt, expectedCaKey, expectedCaPassword, expectedCnames, expectedCountry,
		expectedState, expectedCity, expectedOrg, expectedOrgUnit, expectedSslEmail)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected successful ContainerConfigGenerate call", err != nil)
	test_utils.AssertTrue(t, "ContainerConfigGenerate error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "Result data should be nil", result == nil)
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

			test_utils.AssertEquals(globalT, "ProxyName doesn't match", expectedProxyName, data.ProxyName)
			test_utils.AssertEquals(globalT, "ProxyPort doesn't match", expectedProxyPort, data.ProxyPort)
			test_utils.AssertEquals(globalT, "Server doesn't match", expectedServer, data.Server)
			test_utils.AssertEquals(globalT, "MaxCache doesn't match", expectedMaxCache, data.MaxCache)
			test_utils.AssertEquals(globalT, "Email doesn't match", expectedEmail, data.Email)

			test_utils.AssertEquals(globalT, "CaCertificate doesn't match", expectedCaCrt, data.CaCrt)
			test_utils.AssertEquals(globalT, "CaKey doesn't match", expectedCaKey, data.CaKey)
			test_utils.AssertEquals(globalT, "CaPassword doesn't match", expectedCaPassword, data.CaPassword)
			test_utils.AssertEquals(globalT, "Cnames don't match", fmt.Sprintf("%v", expectedCnames), fmt.Sprintf("%v", data.Cnames))
			test_utils.AssertEquals(globalT, "Country doesn't match", expectedCountry, data.Country)
			test_utils.AssertEquals(globalT, "State doesn't match", expectedState, data.State)
			test_utils.AssertEquals(globalT, "City doesn't match", expectedCity, data.City)
			test_utils.AssertEquals(globalT, "Org doesn't match", expectedOrg, data.Org)
			test_utils.AssertEquals(globalT, "OrgUnit doesn't match", expectedOrgUnit, data.OrgUnit)
			test_utils.AssertEquals(globalT, "SslEmail doesn't match", expectedSslEmail, data.SslEmail)

			// mock response
			return test_utils.GetResponse(200, `{"success": true, "result": [1, 2, 3, 4, 5]}`)
		},
	}
	// Execute
	result, err := proxy.ContainerConfigGenerate(client, expectedProxyName, expectedProxyPort,
		expectedServer, expectedMaxCache, expectedEmail,
		expectedCaCrt, expectedCaKey, expectedCaPassword, expectedCnames, expectedCountry,
		expectedState, expectedCity, expectedOrg, expectedOrgUnit, expectedSslEmail)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected error executing ContainerConfigGenerate", err == nil)
	test_utils.AssertTrue(t, "Result should not be empty", result != nil)
	test_utils.AssertEquals(t, "Result configuration binary doesn't match", fmt.Sprintf("%v", expectedResponseData), fmt.Sprintf("%v", *result))
}
