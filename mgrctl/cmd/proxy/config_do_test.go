// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/mocks"
	proxyApi "github.com/uyuni-project/uyuni-tools/shared/api/proxy"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// common connection details (for generating the client).
var connectionDetails = api.ConnectionDetails{User: "testUser", Password: "testPwd", Server: "testServer"}

// dummy file contents.
const (
	dummyCaCrtContents           = "caCrt contents"
	dummyCaKeyContents           = "caKey contents"
	dummyCaPasswordContents      = "caPwd"
	dummyProxyCrtContents        = "proxyCrt contents"
	dummyProxyKeyContents        = "dummy proxyKey"
	dummyIntermediateCA1Contents = "dummy IntermediateCA 1 contents"
	dummyIntermediateCA2Contents = "dummy IntermediateCA 2 contents"
)

type TestFilePaths struct {
	OutputFilePath          string
	CaCrtFilePath           string
	CaKeyFilePath           string
	ProxyCrtFilePath        string
	ProxyKeyFilePath        string
	IntermediateCA1FilePath string
	IntermediateCA2FilePath string
}

// Helper function to mock a successful API call to login.
func mockSuccessfulLoginAPICall() func(conn *api.ConnectionDetails) (*api.APIClient, error) {
	return func(conn *api.ConnectionDetails) (*api.APIClient, error) {
		client, _ := api.Init(conn)
		client.Client = &mocks.MockClient{
			DoFunc: testutils.SuccessfulLoginTestDo,
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
	intermediateCA1FilePath := createTestFile(testDir, "intermediateCa1.pem", dummyIntermediateCA1Contents, t)
	intermediateCA2FilePath := createTestFile(testDir, "intermediateCa2.pem", dummyIntermediateCA2Contents, t)

	// Return all file paths in a struct
	return TestFilePaths{
		OutputFilePath:          outputFilePath,
		CaCrtFilePath:           caCrtFilePath,
		CaKeyFilePath:           caKeyFilePath,
		ProxyCrtFilePath:        proxyCrtFilePath,
		ProxyKeyFilePath:        proxyKeyFilePath,
		IntermediateCA1FilePath: intermediateCA1FilePath,
		IntermediateCA2FilePath: intermediateCA2FilePath,
	}
}

// tests a failure proxy create config generate command when no connection details are provided.
func TestFailProxyCreateConfigWhenNoConnectionDetailsAreProvided(t *testing.T) {
	// Setup
	testDir := t.TempDir()

	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	flags := &proxyCreateConfigFlags{}
	expectedErrorMessage := "server URL is not provided"

	// Execute
	err := proxyCreateConfig(flags, api.Init, proxyApi.ContainerConfig, proxyApi.ContainerConfigGenerate)

	// Assertions
	testutils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	testutils.AssertTrue(t, "ProxyCreateConfig error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	testutils.AssertTrue(t, "File configuration file stored", !utils.FileExists(expectedOutputFilePath))
}

// tests a failure proxy create config generate command when login fails.
func TestFailProxyCreateConfigWhenLoginFails(t *testing.T) {
	// Setup structures and expected values
	testDir := t.TempDir()

	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	flags := &proxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
	}
	expectedErrorMessage := "Either the password or username is incorrect."

	// Mock login api call to fail
	mockAPIFunc := func(conn *api.ConnectionDetails) (*api.APIClient, error) {
		client, _ := api.Init(conn)
		client.Client = &mocks.MockClient{
			DoFunc: testutils.FailedLoginTestDo,
		}
		return client, nil
	}

	// Execute
	err := proxyCreateConfig(flags, mockAPIFunc, proxyApi.ContainerConfig, proxyApi.ContainerConfigGenerate)

	// Assertions
	testutils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	testutils.AssertTrue(t, "ProxyCreateConfig error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	testutils.AssertTrue(t, "File configuration file stored", !utils.FileExists(expectedOutputFilePath))
}

// tests a failure proxy create config generate command when ProxyCrt is provided but ProxyKey is missing.
func TestFailProxyCreateConfigWhenProxyCrtIsProvidedButProxyKeyIsMissing(t *testing.T) {
	// Setup
	testDir := t.TempDir()

	testFiles := setupTestFiles(t, testDir)
	flags := &proxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
		SSL: proxyConfigSSLFlags{
			Ca: caFlags{
				SSLPair: types.SSLPair{Cert: testFiles.ProxyCrtFilePath},
			},
			Proxy: types.SSLPair{Cert: testFiles.ProxyCrtFilePath},
		},
	}
	expectedErrorMessage := "flag proxyKey is required when flag proxyCrt is provided"

	// Execute
	err := proxyCreateConfig(flags, mockSuccessfulLoginAPICall(), nil, nil)

	// Assertions
	testutils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	testutils.AssertTrue(t, "ProxyCreateConfig error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	testutils.AssertTrue(t, "File configuration file stored", !utils.FileExists(testFiles.OutputFilePath))
}

// tests a failure proxy create config command when proxy config request returns an error.
func TestFailProxyCreateConfigWhenProxyConfigApiRequestFails(t *testing.T) {
	// Setup
	testDir := t.TempDir()

	testFiles := setupTestFiles(t, testDir)
	mockContainerConfigflags := &proxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
		SSL: proxyConfigSSLFlags{
			Ca: caFlags{
				SSLPair: types.SSLPair{Cert: testFiles.CaCrtFilePath},
			},
			Proxy: types.SSLPair{
				Cert: testFiles.ProxyCrtFilePath,
				Key:  testFiles.ProxyKeyFilePath,
			},
		},
	}
	mockContainerConfigGenerateflags := &proxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
		SSL: proxyConfigSSLFlags{
			Ca: caFlags{
				SSLPair: types.SSLPair{
					Cert: testFiles.CaCrtFilePath,
					Key:  testFiles.CaKeyFilePath,
				},
				Password: dummyCaPasswordContents,
			},
		},
	}
	expectedReturnMessage := "Totally unexpected error"

	// Mock containerConfig api calls
	mockContainerConfig := func(_ *api.APIClient, _ proxyApi.ProxyConfigRequest) (*[]int8, error) {
		return nil, errors.New(expectedReturnMessage)
	}
	mockCreateConfigGenerate := func(_ *api.APIClient, _ proxyApi.ProxyConfigGenerateRequest) (*[]int8, error) {
		return nil, errors.New(expectedReturnMessage)
	}

	// Execute providing certs
	err := proxyCreateConfig(
		mockContainerConfigflags, mockSuccessfulLoginAPICall(), mockContainerConfig, mockCreateConfigGenerate,
	)

	// Assertions providing certs call
	testutils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	testutils.AssertTrue(t, "API proxy config return message", strings.HasSuffix(err.Error(), expectedReturnMessage))
	testutils.AssertTrue(t, "File configuration file stored", !utils.FileExists(testFiles.OutputFilePath))

	// Execute generate certs
	err = proxyCreateConfig(
		mockContainerConfigGenerateflags, mockSuccessfulLoginAPICall(), mockContainerConfig, mockCreateConfigGenerate,
	)

	// Assertions generate certs call
	testutils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	testutils.AssertTrue(t, "API proxy config return message", strings.HasSuffix(err.Error(), expectedReturnMessage))
	testutils.AssertTrue(t, "File configuration file stored", !utils.FileExists(testFiles.OutputFilePath))
}

// tests a successful proxy create config command when all parameters provided.
func TestSuccessProxyCreateConfigWhenAllParamsProvidedSuccess(t *testing.T) {
	// Setup
	testDir := t.TempDir()

	testFiles := setupTestFiles(t, testDir)

	output := path.Join(testDir, t.Name())
	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	expectedConfigFileData := []int8{72, 105, 32, 77, 97, 114, 107, 33}

	flags := &proxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
		Proxy: proxyFlags{
			Name:     "testProxy",
			Port:     8080,
			Parent:   "testServer",
			MaxCache: 2048,
			Email:    "example@email.com",
		},
		Output: output,
		SSL: proxyConfigSSLFlags{
			Ca: caFlags{
				SSLPair:      types.SSLPair{Cert: testFiles.CaCrtFilePath},
				Intermediate: []string{testFiles.IntermediateCA1FilePath, testFiles.IntermediateCA2FilePath},
			},
			Proxy: types.SSLPair{
				Cert: testFiles.ProxyCrtFilePath,
				Key:  testFiles.ProxyKeyFilePath,
			},
		},
	}

	// Mock containerConfig api call
	mockContainerConfig := func(_ *api.APIClient, request proxyApi.ProxyConfigRequest) (*[]int8, error) {
		testutils.AssertEquals(t, "Unexpected proxyName", flags.Proxy.Name, request.ProxyName)
		testutils.AssertEquals(t, "Unexpected proxyPort", flags.Proxy.Port, request.ProxyPort)
		testutils.AssertEquals(t, "Unexpected server", flags.Proxy.Parent, request.Server)
		testutils.AssertEquals(t, "Unexpected maxCache", flags.Proxy.MaxCache, request.MaxCache)
		testutils.AssertEquals(t, "Unexpected email", flags.Proxy.Email, request.Email)
		testutils.AssertEquals(t, "Unexpected caCrt", dummyCaCrtContents, request.RootCA)
		testutils.AssertEquals(t, "Unexpected proxyCrt", dummyProxyCrtContents, request.ProxyCrt)
		testutils.AssertEquals(t, "Unexpected proxyKey", dummyProxyKeyContents, request.ProxyKey)
		testutils.AssertEquals(t, "Number of intermediateCAs", 2, len(request.IntermediateCAs))
		testutils.AssertEquals(t, "Unexpected intermediateCA", dummyIntermediateCA1Contents, request.IntermediateCAs[0])
		testutils.AssertEquals(t, "Unexpected intermediateCA", dummyIntermediateCA2Contents, request.IntermediateCAs[1])
		return &expectedConfigFileData, nil
	}

	// Execute
	err := proxyCreateConfig(flags, mockSuccessfulLoginAPICall(), mockContainerConfig, nil)

	// Assertions
	testutils.AssertTrue(t, "Unexpected error executing ProxyCreateConfig", err == nil)
	testutils.AssertTrue(t, "File configuration file was not stored", utils.FileExists(expectedOutputFilePath))

	storedConfigFile := testutils.ReadFileAsBinary(t, expectedOutputFilePath)
	testutils.AssertEquals(t, "File configuration binary doesn't match the response",
		fmt.Sprintf("%v", expectedConfigFileData),
		fmt.Sprintf("%v", storedConfigFile))
}

// tests a successful proxy create config command (with generated certificates) when all parameters provided.
func TestSuccessProxyCreateConfigGenerateWhenAllParamsProvidedSuccess(t *testing.T) {
	// Setup
	testDir := t.TempDir()

	testFiles := setupTestFiles(t, testDir)

	output := path.Join(testDir, t.Name())
	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	expectedConfigFileData := []int8{72, 105, 32, 77, 97, 114, 107, 33}

	flags := &proxyCreateConfigFlags{
		ConnectionDetails: connectionDetails,
		Proxy: proxyFlags{
			Name:     "testProxy",
			Port:     8080,
			Parent:   "testServer",
			MaxCache: 2048,
			Email:    "example@email.com",
		},
		Output: output,
		SSL: proxyConfigSSLFlags{
			SSLCertGenerationFlags: types.SSLCertGenerationFlags{
				Cnames:  []string{"altNameA.example.com", "altNameB.example.com"},
				Country: "testCountry",
				State:   "exampleState",
				City:    "exampleCity",
				Org:     "exampleOrg",
				OU:      "exampleOrgUnit",
				Email:   "sslEmail@example.com",
			},
			Ca: caFlags{
				SSLPair: types.SSLPair{
					Cert: testFiles.CaCrtFilePath,
					Key:  testFiles.CaKeyFilePath,
				},
				Password: dummyCaPasswordContents,
			},
		},
	}

	// Mock api client & containerConfig
	mockCreateConfigGenerate := func(_ *api.APIClient, request proxyApi.ProxyConfigGenerateRequest) (*[]int8, error) {
		testutils.AssertEquals(t, "Unexpected proxyName", flags.Proxy.Name, request.ProxyName)
		testutils.AssertEquals(t, "Unexpected proxyPort", flags.Proxy.Port, request.ProxyPort)
		testutils.AssertEquals(t, "Unexpected server", flags.Proxy.Parent, request.Server)
		testutils.AssertEquals(t, "Unexpected maxCache", flags.Proxy.MaxCache, request.MaxCache)
		testutils.AssertEquals(t, "Unexpected email", flags.Proxy.Email, request.Email)
		testutils.AssertEquals(t, "Unexpected caCrt", dummyCaCrtContents, request.CaCrt)
		testutils.AssertEquals(t, "Unexpected caKey", dummyCaKeyContents, request.CaKey)
		testutils.AssertEquals(t, "Unexpected caPassword", dummyCaPasswordContents, request.CaPassword)
		testutils.AssertEquals(t, "Unexpected cnames", fmt.Sprintf("%v", flags.SSL.Cnames), fmt.Sprintf("%v", request.Cnames))
		testutils.AssertEquals(t, "Unexpected country", flags.SSL.Country, request.Country)
		testutils.AssertEquals(t, "Unexpected state", flags.SSL.State, request.State)
		testutils.AssertEquals(t, "Unexpected city", flags.SSL.City, request.City)
		testutils.AssertEquals(t, "Unexpected org", flags.SSL.Org, request.Org)
		testutils.AssertEquals(t, "Unexpected orgUnit", flags.SSL.OU, request.OrgUnit)
		testutils.AssertEquals(t, "Unexpected sslEmail", flags.SSL.Email, request.SSLEmail)
		return &expectedConfigFileData, nil
	}

	// Execute
	err := proxyCreateConfig(flags, mockSuccessfulLoginAPICall(), nil, mockCreateConfigGenerate)

	// Assertions
	testutils.AssertTrue(t, "Unexpected error executing ProxyCreateConfigGenerate", err == nil)
	testutils.AssertTrue(t, "File configuration file was not stored", utils.FileExists(expectedOutputFilePath))

	storedConfigFile := testutils.ReadFileAsBinary(t, expectedOutputFilePath)
	testutils.AssertEquals(t, "File configuration binary doesn't match the response",
		fmt.Sprintf("%v", expectedConfigFileData),
		fmt.Sprintf("%v", storedConfigFile))
}
