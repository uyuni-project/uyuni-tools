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

// dummy file contents.
const dummyCaCrtContents = "caCrt contents"
const dummyCaKeyContents = "caKey contents"
const dummyCaPasswordContents = "caPwd"

// tests a failure proxy create config generate command when no connection details are provided.
func TestFailProxyCreateConfigGenerateWhenNoConnectionDetailsAreProvided(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	flags := &proxy.ProxyCreateConfigGenerateFlags{}

	expectedErrorMessage := "server URL is not provided"

	// Execute
	err := proxy.ProxyCreateConfigGenerate(flags, api.Init, proxyApi.ContainerConfigGenerate)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfigGenerate", err != nil)
	test_utils.AssertTrue(t, "ProxyCreateConfigGenerate error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(expectedOutputFilePath))
}

// tests a failure proxy create config generate command when login fails.
func TestFailProxyCreateConfigGenerateWhenLoginFails(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	flags := &proxy.ProxyCreateConfigGenerateFlags{
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
	err := proxy.ProxyCreateConfigGenerate(flags, mockAPIFunc, proxyApi.ContainerConfigGenerate)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfigGenerate", err != nil)
	test_utils.AssertTrue(t, "ProxyCreateConfigGenerate error message", strings.HasSuffix(err.Error(), expectedErrorMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(expectedOutputFilePath))
}

// tests a failure proxy create config generate command when proxy config request returns an error.
func TestFailProxyCreateConfigGenerateWhenProxyConfigApiRequestFails(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	expectedReturnMessage := "Totally unexpected error"

	caCrtFilePath := createTestFile(testDir, "ca.pem", dummyCaCrtContents, t)
	caKeyFilePath := createTestFile(testDir, "caKey.pem", dummyCaKeyContents, t)
	caPwdFilePath := createTestFile(testDir, "pass.txt", dummyCaPasswordContents, t)

	flags := &proxy.ProxyCreateConfigGenerateFlags{
		ProxyCreateConfigBaseFlags: proxy.ProxyCreateConfigBaseFlags{
			ProxyName:         "testProxy",
			Server:            "testServer",
			Email:             "example@email.com",
			ConnectionDetails: connectionDetails,
		},
		CaCrt:      caCrtFilePath,
		CaKey:      caKeyFilePath,
		CaPassword: caPwdFilePath,
	}

	// Mock api client & containerConfig
	mockAPIFunc := func(conn *api.ConnectionDetails) (*api.APIClient, error) {
		client, _ := api.Init(conn)
		client.Client = &mocks.MockClient{
			DoFunc: test_utils.SuccessfulLoginTestDo,
		}
		return client, nil
	}
	mockCreateConfigGenerate := func(client *api.APIClient, proxyName string, proxyPort int, server string, maxCache int, email string, caCrt string, caKey string, caPassword string, cnames []string, country string, state string, city string, org string, orgUnit string, sslEmail string) (*[]int8, error) {
		return nil, errors.New(expectedReturnMessage)
	}

	// Execute
	err := proxy.ProxyCreateConfigGenerate(flags, mockAPIFunc, mockCreateConfigGenerate)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected success execution of ProxyCreateConfigGenerate", err != nil)
	test_utils.AssertTrue(t, "API proxy config return message", strings.HasSuffix(err.Error(), expectedReturnMessage))
	test_utils.AssertTrue(t, "File configuration file stored", !utils.FileExists(expectedOutputFilePath))
}

// tests a successful proxy create config generate command with all parameters provided.
func TestProxyCreateConfigGenerateWhenAllParamsProvidedSuccess(t *testing.T) {
	// Setup
	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	caCrtFilePath := createTestFile(testDir, "ca.pem", dummyCaCrtContents, t)
	caKeyFilePath := createTestFile(testDir, "caKey.pem", dummyCaKeyContents, t)
	caPwdFilePath := createTestFile(testDir, "pass.txt", dummyCaPasswordContents, t)

	output := path.Join(testDir, t.Name())
	expectedOutputFilePath := path.Join(testDir, t.Name()+".tar.gz")
	expectedConfigFileData := []int8{72, 105, 32, 77, 97, 114, 107, 33}

	flags := &proxy.ProxyCreateConfigGenerateFlags{
		ProxyCreateConfigBaseFlags: proxy.ProxyCreateConfigBaseFlags{
			ProxyName:         "testProxy",
			ProxyPort:         8080,
			Server:            "testServer",
			MaxCache:          2048,
			Email:             "example@email.com",
			Output:            output,
			ConnectionDetails: connectionDetails,
		},
		CaCrt:      caCrtFilePath,
		CaKey:      caKeyFilePath,
		CaPassword: caPwdFilePath,
		CNames:     []string{"altNameA.example.com", "altNameB.example.com"},
		Country:    "testCountry",
		State:      "exampleState",
		City:       "exampleCity",
		Org:        "exampleOrg",
		OrgUnit:    "exampleOrgUnit",
		SslEmail:   "sslEmail@example.com",
	}

	// Mock api client & containerConfig
	mockAPIFunc := func(conn *api.ConnectionDetails) (*api.APIClient, error) {
		client, _ := api.Init(conn)
		client.Client = &mocks.MockClient{
			DoFunc: test_utils.SuccessfulLoginTestDo,
		}
		return client, nil
	}
	mockCreateConfigGenerate := func(client *api.APIClient, proxyName string, proxyPort int, server string, maxCache int, email string, caCrt string, caKey string, caPassword string, cnames []string, country string, state string, city string, org string, orgUnit string, sslEmail string) (*[]int8, error) {
		test_utils.AssertEquals(t, "Unexpected proxyName", flags.ProxyName, proxyName)
		test_utils.AssertEquals(t, "Unexpected proxyPort", flags.ProxyPort, proxyPort)
		test_utils.AssertEquals(t, "Unexpected server", flags.Server, server)
		test_utils.AssertEquals(t, "Unexpected maxCache", flags.MaxCache, maxCache)
		test_utils.AssertEquals(t, "Unexpected email", flags.Email, email)
		test_utils.AssertEquals(t, "Unexpected caCrt", dummyCaCrtContents, caCrt)
		test_utils.AssertEquals(t, "Unexpected caKey", dummyCaKeyContents, caKey)
		test_utils.AssertEquals(t, "Unexpected caPassword", dummyCaPasswordContents, caPassword)
		test_utils.AssertEquals(t, "Unexpected cnames", fmt.Sprintf("%v", flags.CNames), fmt.Sprintf("%v", cnames))
		test_utils.AssertEquals(t, "Unexpected country", flags.Country, country)
		test_utils.AssertEquals(t, "Unexpected state", flags.State, state)
		test_utils.AssertEquals(t, "Unexpected city", flags.City, city)
		test_utils.AssertEquals(t, "Unexpected org", flags.Org, org)
		test_utils.AssertEquals(t, "Unexpected orgUnit", flags.OrgUnit, orgUnit)
		test_utils.AssertEquals(t, "Unexpected sslEmail", flags.SslEmail, sslEmail)
		return &expectedConfigFileData, nil
	}

	// Execute
	err := proxy.ProxyCreateConfigGenerate(flags, mockAPIFunc, mockCreateConfigGenerate)

	// Assertions
	test_utils.AssertTrue(t, "Unexpected error executing ProxyCreateConfigGenerate", err == nil)
	test_utils.AssertTrue(t, "File configuration file was not stored", utils.FileExists(expectedOutputFilePath))

	storedConfigFile := test_utils.ReadFileAsBinary(t, expectedOutputFilePath)
	test_utils.AssertEquals(t, "File configuration binary doesn't match the response", fmt.Sprintf("%v", expectedConfigFileData), fmt.Sprintf("%v",
		storedConfigFile))
}
