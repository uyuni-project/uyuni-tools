// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy_test

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/proxy"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/mocks"
	proxyApi "github.com/uyuni-project/uyuni-tools/shared/api/proxy"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// common connection details (for generating the client).
const user = "testUser"
const password = "testPwd"
const server = "testServer"

var connectionDetails = api.ConnectionDetails{User: user, Password: password, Server: server}

// dummy file contents.
const dummyCaCrtContents = "caCrt contents"
const dummyCaKeyContents = "caKey contents"
const dummyCaPasswordContents = "caPwd"
const dummyProxyCrtContents = "proxyCrt contents"
const dummyProxyKeyContents = "dummy proxyKey"
const dummyIntermediateCA1Contents = "dummy IntermediateCA 1 contents"
const dummyIntermediateCA2Contents = "dummy IntermediateCA 2 contents"

type TestFilePaths struct {
	OutputFilePath          string
	CaCrtFilePath           string
	CaKeyFilePath           string
	ProxyCrtFilePath        string
	ProxyKeyFilePath        string
	CaPwdFilePath           string
	IntermediateCA1FilePath string
	IntermediateCA2FilePath string
}

// Helper function to mock a successful API call to login.
func mockSuccessfulLoginApiCall() func(conn *api.ConnectionDetails) (*api.APIClient, error) {
	return func(conn *api.ConnectionDetails) (*api.APIClient, error) {
		client, _ := api.Init(conn)
		client.Client = &mocks.MockClient{
			DoFunc: test_utils.SuccessfulLoginTestDo,
		}
		return client, nil
	}
}

// Helper function to create a test files with given contents.
func setupTestFiles(t *testing.T, testDir string) TestFilePaths {
	outputFilePath := path.Join(testDir, t.Name()+".tar.gz")

	caCrtFilePath := createTestFile(testDir, "ca.pem", dummyCaCrtContents, t)
	caKeyFilePath := createTestFile(testDir, "caKey.pem", dummyCaKeyContents, t)
	proxyCrtFilePath := createTestFile(testDir, "proxyCrt.pem", dummyProxyCrtContents, t)
	proxyKeyFilePath := createTestFile(testDir, "proxyKey.txt", dummyProxyKeyContents, t)
	caPwdFilePath := createTestFile(testDir, "pass.txt", dummyCaPasswordContents, t)
	intermediateCA1FilePath := createTestFile(testDir, "intermediateCa1.pem", dummyIntermediateCA1Contents, t)
	intermediateCA2FilePath := createTestFile(testDir, "intermediateCa2.pem", dummyIntermediateCA2Contents, t)

	// Return all file paths in a struct
	return TestFilePaths{
		OutputFilePath:          outputFilePath,
		CaCrtFilePath:           caCrtFilePath,
		CaKeyFilePath:           caKeyFilePath,
		ProxyCrtFilePath:        proxyCrtFilePath,
		ProxyKeyFilePath:        proxyKeyFilePath,
		CaPwdFilePath:           caPwdFilePath,
		IntermediateCA1FilePath: intermediateCA1FilePath,
		IntermediateCA2FilePath: intermediateCA2FilePath,
	}
}

// tests a failure proxy create config generate command when no connection details are provided.
func TestFailProxyCreateConfigWhenNoConnectionDetailsAreProvided(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	flags := &proxy.ProxyCreateConfigFlags{}
	expectedErrorMessage := "server URL is not provided"

	// Execute
	err := proxy.ProxyCreateConfig(flags, api.Init, proxyApi.ContainerConfig, proxyApi.ContainerConfigGenerate)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	test_utils.AssertTrue(t, "ProxyCreateConfig error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(expectedOutputFilePath))
}

// tests a failure proxy create config generate command when login fails.
func TestFailProxyCreateConfigWhenLoginFails(t *testing.T) {
	// Setup structures and expected values
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	flags := &proxy.ProxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
	}
	expectedErrorMessage := "Either the password or username is incorrect."

	// Mock login api call to fail
	mockAPIFunc := func(conn *api.ConnectionDetails) (*api.APIClient, error) {
		client, _ := api.Init(conn)
		client.Client = &mocks.MockClient{
			DoFunc: test_utils.FailedLoginTestDo,
		}
		return client, nil
	}

	// Execute
	err := proxy.ProxyCreateConfig(flags, mockAPIFunc, proxyApi.ContainerConfig, proxyApi.ContainerConfigGenerate)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	test_utils.AssertTrue(t, "ProxyCreateConfig error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(expectedOutputFilePath))
}

// tests a failure proxy create config generate command when ProxyCrt is provided but ProxyKey is missing.
func TestFailProxyCreateConfigWhenProxyCrtIsProvidedButProxyKeyIsMissing(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	testFiles := setupTestFiles(t, testDir)
	flags := &proxy.ProxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
		CaCrt:             testFiles.CaCrtFilePath,
		ProxyCrt:          testFiles.ProxyCrtFilePath,
	}
	expectedErrorMessage := "flag proxyKey is required when flag proxyCrt is provided"

	// Execute
	err := proxy.ProxyCreateConfig(flags, mockSuccessfulLoginApiCall(), nil, nil)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	test_utils.AssertTrue(t, "ProxyCreateConfig error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(testFiles.OutputFilePath))
}

// tests a failure proxy create config command when proxy config request returns an error.
func TestFailProxyCreateConfigWhenProxyConfigApiRequestFails(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	testFiles := setupTestFiles(t, testDir)
	mockContainerConfigflags := &proxy.ProxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
		CaCrt:             testFiles.CaCrtFilePath,
		ProxyCrt:          testFiles.ProxyCrtFilePath,
		ProxyKey:          testFiles.ProxyKeyFilePath,
	}
	mockContainerConfigGenerateflags := &proxy.ProxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
		CaCrt:             testFiles.CaCrtFilePath,
		CaKey:             testFiles.CaKeyFilePath,
		CaPassword:        testFiles.CaPwdFilePath,
	}
	expectedReturnMessage := "Totally unexpected error"

	// Mock containerConfig api calls
	mockContainerConfig := func(client *api.APIClient, request proxyApi.ProxyConfigRequest) (*[]int8, error) {
		return nil, errors.New(expectedReturnMessage)
	}
	mockCreateConfigGenerate := func(client *api.APIClient, request proxyApi.ProxyConfigGenerateRequest) (*[]int8, error) {
		return nil, errors.New(expectedReturnMessage)
	}

	// Execute providing certs
	err := proxy.ProxyCreateConfig(mockContainerConfigflags, mockSuccessfulLoginApiCall(), mockContainerConfig, mockCreateConfigGenerate)

	// Assertions providing certs call
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	test_utils.AssertTrue(t, "API proxy config return message", strings.HasSuffix(err.Error(), expectedReturnMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(testFiles.OutputFilePath))

	// Execute generate certs
	err = proxy.ProxyCreateConfig(mockContainerConfigGenerateflags, mockSuccessfulLoginApiCall(), mockContainerConfig, mockCreateConfigGenerate)

	// Assertions generate certs call
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	test_utils.AssertTrue(t, "API proxy config return message", strings.HasSuffix(err.Error(), expectedReturnMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(testFiles.OutputFilePath))
}

// tests a successful proxy create config command when all parameters provided.
func TestSuccessProxyCreateConfigWhenAllParamsProvidedSuccess(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	testFiles := setupTestFiles(t, testDir)

	output := path.Join(testDir, t.Name())
	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	expectedConfigFileData := []int8{72, 105, 32, 77, 97, 114, 107, 33}

	flags := &proxy.ProxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
		ProxyName:         "testProxy",
		ProxyPort:         8080,
		Server:            "testServer",
		MaxCache:          2048,
		Email:             "example@email.com",
		Output:            output,
		CaCrt:             testFiles.CaCrtFilePath,
		ProxyCrt:          testFiles.ProxyCrtFilePath,
		ProxyKey:          testFiles.ProxyKeyFilePath,
		IntermediateCAs:   []string{testFiles.IntermediateCA1FilePath, testFiles.IntermediateCA2FilePath},
	}

	// Mock containerConfig api call
	mockContainerConfig := func(client *api.APIClient, request proxyApi.ProxyConfigRequest) (*[]int8, error) {
		test_utils.AssertEquals(t, "Unexpected proxyName", flags.ProxyName, request.ProxyName)
		test_utils.AssertEquals(t, "Unexpected proxyPort", flags.ProxyPort, request.ProxyPort)
		test_utils.AssertEquals(t, "Unexpected server", flags.Server, request.Server)
		test_utils.AssertEquals(t, "Unexpected maxCache", flags.MaxCache, request.MaxCache)
		test_utils.AssertEquals(t, "Unexpected email", flags.Email, request.Email)
		test_utils.AssertEquals(t, "Unexpected caCrt", dummyCaCrtContents, request.RootCA)
		test_utils.AssertEquals(t, "Unexpected proxyCrt", dummyProxyCrtContents, request.ProxyCrt)
		test_utils.AssertEquals(t, "Unexpected proxyKey", dummyProxyKeyContents, request.ProxyKey)
		test_utils.AssertEquals(t, "Number of intermediateCAs", 2, len(request.IntermediateCAs))
		test_utils.AssertEquals(t, "Unexpected intermediateCA", dummyIntermediateCA1Contents, request.IntermediateCAs[0])
		test_utils.AssertEquals(t, "Unexpected intermediateCA", dummyIntermediateCA2Contents, request.IntermediateCAs[1])
		return &expectedConfigFileData, nil
	}

	// Execute
	err := proxy.ProxyCreateConfig(flags, mockSuccessfulLoginApiCall(), mockContainerConfig, nil)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected error executing ProxyCreateConfig", err == nil)
	test_utils.AssertTrue(t, "File configuration file was not stored", utils.FileExists(expectedOutputFilePath))

	storedConfigFile := test_utils.ReadFileAsBinary(t, expectedOutputFilePath)
	test_utils.AssertEquals(t, "File configuration binary doesn't match the response",
		fmt.Sprintf("%v", expectedConfigFileData),
		fmt.Sprintf("%v", storedConfigFile))
}

// tests a successful proxy create config command (with generated certificates) when all parameters provided.
func TestSuccessProxyCreateConfigGenerateWhenAllParamsProvidedSuccess(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	testFiles := setupTestFiles(t, testDir)

	output := path.Join(testDir, t.Name())
	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	expectedConfigFileData := []int8{72, 105, 32, 77, 97, 114, 107, 33}

	flags := &proxy.ProxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
		ProxyName:         "testProxy",
		ProxyPort:         8080,
		Server:            "testServer",
		MaxCache:          2048,
		Email:             "example@email.com",
		Output:            output,
		CaCrt:             testFiles.CaCrtFilePath,
		CaKey:             testFiles.CaKeyFilePath,
		CaPassword:        testFiles.CaPwdFilePath,
		CNames:            []string{"altNameA.example.com", "altNameB.example.com"},
		Country:           "testCountry",
		State:             "exampleState",
		City:              "exampleCity",
		Org:               "exampleOrg",
		OrgUnit:           "exampleOrgUnit",
		SslEmail:          "sslEmail@example.com",
	}

	// Mock api client & containerConfig
	mockCreateConfigGenerate := func(client *api.APIClient, request proxyApi.ProxyConfigGenerateRequest) (*[]int8, error) {
		test_utils.AssertEquals(t, "Unexpected proxyName", flags.ProxyName, request.ProxyName)
		test_utils.AssertEquals(t, "Unexpected proxyPort", flags.ProxyPort, request.ProxyPort)
		test_utils.AssertEquals(t, "Unexpected server", flags.Server, request.Server)
		test_utils.AssertEquals(t, "Unexpected maxCache", flags.MaxCache, request.MaxCache)
		test_utils.AssertEquals(t, "Unexpected email", flags.Email, request.Email)
		test_utils.AssertEquals(t, "Unexpected caCrt", dummyCaCrtContents, request.CaCrt)
		test_utils.AssertEquals(t, "Unexpected caKey", dummyCaKeyContents, request.CaKey)
		test_utils.AssertEquals(t, "Unexpected caPassword", dummyCaPasswordContents, request.CaPassword)
		test_utils.AssertEquals(t, "Unexpected cnames", fmt.Sprintf("%v", flags.CNames), fmt.Sprintf("%v", request.Cnames))
		test_utils.AssertEquals(t, "Unexpected country", flags.Country, request.Country)
		test_utils.AssertEquals(t, "Unexpected state", flags.State, request.State)
		test_utils.AssertEquals(t, "Unexpected city", flags.City, request.City)
		test_utils.AssertEquals(t, "Unexpected org", flags.Org, request.Org)
		test_utils.AssertEquals(t, "Unexpected orgUnit", flags.OrgUnit, request.OrgUnit)
		test_utils.AssertEquals(t, "Unexpected sslEmail", flags.SslEmail, request.SslEmail)
		return &expectedConfigFileData, nil
	}

	// Execute
	err := proxy.ProxyCreateConfig(flags, mockSuccessfulLoginApiCall(), nil, mockCreateConfigGenerate)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected error executing ProxyCreateConfigGenerate", err == nil)
	test_utils.AssertTrue(t, "File configuration file was not stored", utils.FileExists(expectedOutputFilePath))

	storedConfigFile := test_utils.ReadFileAsBinary(t, expectedOutputFilePath)
	test_utils.AssertEquals(t, "File configuration binary doesn't match the response",
		fmt.Sprintf("%v", expectedConfigFileData),
		fmt.Sprintf("%v", storedConfigFile))
}
