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
const dummyRootCAContents = "rootCA contents"
const dummyProxyCrtContents = "proxyCrt contents"
const dummyProxyKeyContents = "dummy proxyKey"
const dummyIntermediateCA1Contents = "dummy IntermediateCA 1 contents"
const dummyIntermediateCA2Contents = "dummy IntermediateCA 2 contents"

// tests a failure proxy create config generate command when no connection details are provided.
func TestFailProxyCreateConfigWhenNoConnectionDetailsAreProvided(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	flags := &proxy.ProxyCreateConfigFlags{}

	expectedErrorMessage := "server URL is not provided"

	// Execute
	err := proxy.ProxyCreateConfig(flags, api.Init, proxyApi.ContainerConfig)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	test_utils.AssertTrue(t, "ProxyCreateConfig error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(expectedOutputFilePath))
}

// tests a failure proxy create config generate command when login fails.
func TestFailProxyCreateConfigWhenLoginFails(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	flags := &proxy.ProxyCreateConfigFlags{
		ProxyCreateConfigBaseFlags: proxy.ProxyCreateConfigBaseFlags{
			ConnectionDetails: connectionDetails,
		},
	}
	expectedErrorMessage := "Either the password or username is incorrect."

	// Mock api client & containerConfig
	mockAPIFunc := func(conn *api.ConnectionDetails) (*api.APIClient, error) {
		client, _ := api.Init(conn)
		client.Client = &mocks.MockClient{
			DoFunc: test_utils.FailedLoginTestDo,
		}
		return client, nil
	}

	// Execute
	err := proxy.ProxyCreateConfig(flags, mockAPIFunc, proxyApi.ContainerConfig)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	test_utils.AssertTrue(t, "ProxyCreateConfig error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(expectedOutputFilePath))
}

// tests a failure proxy create config generate command when proxy config request returns an error.
func TestFailProxyCreateConfigWhenProxyConfigApiRequestFails(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	expectedReturnMessage := "Totally unexpected error"

	rootCaFilePath := createTestFile(testDir, "rootCA.pem", dummyRootCAContents, t)
	proxyCrtFilePath := createTestFile(testDir, "proxyCrt.pem", dummyProxyCrtContents, t)
	proxyKeyFilePath := createTestFile(testDir, "proxyKey.txt", dummyProxyKeyContents, t)

	flags := &proxy.ProxyCreateConfigFlags{
		ProxyCreateConfigBaseFlags: proxy.ProxyCreateConfigBaseFlags{
			ProxyName:         "testProxy",
			Server:            "testServer",
			Email:             "example@email.com",
			ConnectionDetails: connectionDetails,
		},
		RootCA:   rootCaFilePath,
		ProxyCrt: proxyCrtFilePath,
		ProxyKey: proxyKeyFilePath,
	}

	// Mock api client & containerConfig
	mockAPIFunc := func(conn *api.ConnectionDetails) (*api.APIClient, error) {
		client, _ := api.Init(conn)
		client.Client = &mocks.MockClient{
			DoFunc: test_utils.SuccessfulLoginTestDo,
		}
		return client, nil
	}
	mockContainerConfig := func(client *api.APIClient, proxyName string, proxyPort int,
		server string, maxCache int, email string,
		rootCA string, proxyCrt string, proxyKey string, intermediateCAs []string) (*[]int8, error) {
		return nil, errors.New(expectedReturnMessage)
	}

	// Execute
	err := proxy.ProxyCreateConfig(flags, mockAPIFunc, mockContainerConfig)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfig", err != nil)
	test_utils.AssertTrue(t, "API proxy config return message", strings.HasSuffix(err.Error(), expectedReturnMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(expectedOutputFilePath))
}

// tests a successful proxy create config generate command with all parameters provided.
func TestProxyCreateConfigWhenAllParamsProvidedSuccess(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	rootCaFilePath := createTestFile(testDir, "rootCA.pem", dummyRootCAContents, t)
	proxyCrtFilePath := createTestFile(testDir, "proxyCrt.pem", dummyProxyCrtContents, t)
	proxyKeyFilePath := createTestFile(testDir, "proxyKey.txt", dummyProxyKeyContents, t)
	intermediateCA1FilePath := createTestFile(testDir, "intermediateCa1.pem", dummyIntermediateCA1Contents, t)
	intermediateCA2FilePath := createTestFile(testDir, "intermediateCa2.pem", dummyIntermediateCA2Contents, t)

	output := path.Join(testDir, t.Name())
	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	expectedConfigFileData := []int8{72, 105, 32, 77, 97, 114, 107, 33}

	flags := &proxy.ProxyCreateConfigFlags{
		ProxyCreateConfigBaseFlags: proxy.ProxyCreateConfigBaseFlags{
			ProxyName:         "testProxy",
			ProxyPort:         8080,
			Server:            "testServer",
			MaxCache:          2048,
			Email:             "example@email.com",
			Output:            output,
			ConnectionDetails: connectionDetails,
		},
		RootCA:          rootCaFilePath,
		ProxyCrt:        proxyCrtFilePath,
		ProxyKey:        proxyKeyFilePath,
		IntermediateCAs: []string{intermediateCA1FilePath, intermediateCA2FilePath},
	}

	// Mock api client & containerConfig
	mockAPIFunc := func(conn *api.ConnectionDetails) (*api.APIClient, error) {
		client, _ := api.Init(conn)
		client.Client = &mocks.MockClient{
			DoFunc: test_utils.SuccessfulLoginTestDo,
		}
		return client, nil
	}
	mockContainerConfig := func(client *api.APIClient, proxyName string, proxyPort int,
		server string, maxCache int, email string,
		rootCA string, proxyCrt string, proxyKey string, intermediateCAs []string) (*[]int8, error) {
		test_utils.AssertEquals(t, "Unexpected proxyName", flags.ProxyName, proxyName)
		test_utils.AssertEquals(t, "Unexpected proxyPort", flags.ProxyPort, proxyPort)
		test_utils.AssertEquals(t, "Unexpected server", flags.Server, server)
		test_utils.AssertEquals(t, "Unexpected maxCache", flags.MaxCache, maxCache)
		test_utils.AssertEquals(t, "Unexpected email", flags.Email, email)
		test_utils.AssertEquals(t, "Unexpected rootCA", dummyRootCAContents, rootCA)
		test_utils.AssertEquals(t, "Unexpected proxyCrt", dummyProxyCrtContents, proxyCrt)
		test_utils.AssertEquals(t, "Unexpected proxyKey", dummyProxyKeyContents, proxyKey)
		test_utils.AssertEquals(t, "Number of intermediateCAs", 2, len(intermediateCAs))
		test_utils.AssertEquals(t, "Unexpected intermediateCA", dummyIntermediateCA1Contents, intermediateCAs[0])
		test_utils.AssertEquals(t, "Unexpected intermediateCA", dummyIntermediateCA2Contents, intermediateCAs[1])
		return &expectedConfigFileData, nil
	}

	// Execute
	err := proxy.ProxyCreateConfig(flags, mockAPIFunc, mockContainerConfig)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected error executing ProxyCreateConfig", err == nil)
	test_utils.AssertTrue(t, "File configuration file was not stored", utils.FileExists(expectedOutputFilePath))

	storedConfigFile := test_utils.ReadFileAsBinary(t, expectedOutputFilePath)
	test_utils.AssertEquals(t, "File configuration binary doesn't match the response", fmt.Sprintf("%v", expectedConfigFileData), fmt.Sprintf("%v",
		storedConfigFile))
}
