// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"syscall"
	"testing"

	expect "github.com/Netflix/go-expect"
	"github.com/chai2010/gettext-go"
	"github.com/spf13/cobra"
	l10n_utils "github.com/uyuni-project/uyuni-tools/shared/l10n/utils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

type askTestData struct {
	value           string
	expectedMessage string
	min             int
	max             int
	checker         func(string) bool
}

func setupConsole(t *testing.T) (*expect.Console, func()) {
	// Set english locale to not depend on the system one
	gettext.BindLocale(gettext.New("", "", l10n_utils.New("")))
	gettext.SetLanguage("en")

	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		t.Errorf("Failed to create fake console")
	}

	origStdin := syscall.Stdin
	origOsStdin := os.Stdin
	origStdout := os.Stdout

	syscall.Stdin = int(c.Tty().Fd())
	os.Stdin = c.Tty()
	os.Stdout = c.Tty()

	return c, func() {
		syscall.Stdin = origStdin
		os.Stdin = origOsStdin
		os.Stdout = origStdout
		c.Close()
	}
}

func TestAskIfMissing(t *testing.T) {
	c, teardown := setupConsole(t)
	defer teardown()

	fChecker := func(v string) bool {
		if !strings.Contains(v, "f") {
			fmt.Println("Has to contain an 'f'")
			return false
		}
		return true
	}

	data := []askTestData{
		{value: "\n", expectedMessage: "A value is required", min: 1, max: 5},
		{value: "superlong\n", expectedMessage: "Has to be less than 5 characters long", min: 1, max: 5},
		{value: "a\n", expectedMessage: "Has to be more than 2 characters long", min: 2, max: 5},
		{value: "booh\n", expectedMessage: "Has to contain an 'f'", min: 0, max: 0, checker: fChecker},
	}

	for i, testCase := range data {
		go func() {
			sendInput(t, i, c, "Prompted value:", testCase.value, testCase.expectedMessage)
			// Send a good value
			sendInput(t, i, c, "Prompted value:", "foo\n", "")
		}()

		var value string
		AskIfMissing(&value, "Prompted value", testCase.min, testCase.max, testCase.checker)
		if value != "foo" {
			t.Errorf("Testcase %d: Expected 'foo', got '%s' value", i, value)
		}
	}
}

func TestCheckValidPassword(t *testing.T) {
	c, teardown := setupConsole(t)
	defer teardown()

	data := []askTestData{
		{value: "\n", expectedMessage: "A value is required", min: 1, max: 5},
		{value: "superlong\n", expectedMessage: "Has to be less than 5 characters long", min: 1, max: 5},
		{value: "a\n", expectedMessage: "Has to be more than 2 characters long", min: 2, max: 5},
	}

	for i, testCase := range data {
		go func() {
			sendInput(t, i, c, "Prompted password:", testCase.value, testCase.expectedMessage)
			// Send a good password
			sendInput(t, i, c, "Prompted password: ", "foo\n", "")
			sendInput(t, i, c, "Confirm the password: ", "foo\n", "")
		}()

		var value string
		AskPasswordIfMissing(&value, "Prompted password", testCase.min, testCase.max)
		if value != "foo" {
			t.Errorf("Testcase %d: Expected 'foo', got '%s' value", i, value)
		}
	}
}

func TestPasswordMismatch(t *testing.T) {
	c, teardown := setupConsole(t)
	defer teardown()

	go func() {
		sendInput(t, 1, c, "Prompted password: ", "password1\n", "")
		sendInput(t, 1, c, "Confirm the password: ", "password2\n", "")

		if _, err := c.ExpectString("Two different passwords have been provided"); err != nil {
			t.Errorf("Expected message error: %s", err)
		}

		// Send a good password
		sendInput(t, 1, c, "Prompted password: ", "foo\n", "")
		sendInput(t, 1, c, "Confirm the password: ", "foo\n", "")
	}()

	var value string
	AskPasswordIfMissing(&value, "Prompted password", 1, 20)
	if value != "foo" {
		t.Errorf("Expected 'foo', got '%s' value", value)
	}
}

func sendInput(
	t *testing.T,
	testcase int,
	c *expect.Console,
	expectedPrompt string,
	value string,
	expectedMessage string,
) {
	if _, err := c.ExpectString(expectedPrompt); err != nil {
		t.Errorf("Testcase %d: Expected prompt error: %s", testcase, err)
	}
	if _, err := c.Send(value); err != nil {
		t.Errorf("Testcase %d: Failed to send value to fake console: %s", testcase, err)
	}
	t.Logf("Value sent: '%s'", value)
	if expectedMessage == "" {
		return
	}

	if _, err := c.Expect(expect.Regexp(regexp.MustCompile(expectedMessage))); err != nil {
		t.Errorf("Testcase %d: Expected '%s' message: %s", testcase, expectedMessage, err)
	}
	if expectedMessage == "" {
		return
	}
}

func TestComputePTF(t *testing.T) {
	// Constants
	const (
		defaultPtfID      = "27977"
		defaultUser       = "150158"
		defaultSuffix     = "ptf"
		baseRegistryHost  = "registry.suse.com"
		defaultRegistry50 = "registry.suse.com/suse/manager/5.0/x86_64"
		defaultRegistry51 = "registry.suse.com/suse/multi-linux-manager/5.1/x86_64"
	)

	tests := []struct {
		name                 string
		registry             string
		user                 string
		ptfID                string
		fullImage            string
		suffix               string
		expected             string
		expectedErrorMessage string
	}{
		// Success cases - 5.0 Manager
		{
			name:      "success 5.0 container with 5.0 registry",
			registry:  defaultRegistry50,
			fullImage: defaultRegistry50 + "/proxy-tftpd:5.0.0",
			expected:  "registry.suse.com/a/150158/27977/suse/manager/5.0/x86_64/proxy-tftpd:latest-ptf-27977",
		},
		{
			name:      "success 5.0 rpm container with 5.0 registry",
			registry:  defaultRegistry50,
			fullImage: "localhost/suse/manager/5.0/x86_64/proxy-ssh:5.0.0",
			expected:  "registry.suse.com/a/150158/27977/suse/manager/5.0/x86_64/proxy-ssh:latest-ptf-27977",
		},
		{
			name:      "success 5.0 container and base registry host",
			registry:  baseRegistryHost,
			fullImage: defaultRegistry50 + "/proxy-tftpd:latest",
			expected:  "registry.suse.com/a/150158/27977/suse/manager/5.0/x86_64/proxy-tftpd:latest-ptf-27977",
		},
		{
			name:      "success 5.0 container and custom registry",
			registry:  "mysccregistry.com",
			fullImage: defaultRegistry50 + "/proxy-helm:latest",
			expected:  "mysccregistry.com/a/150158/27977/suse/manager/5.0/x86_64/proxy-helm:latest-ptf-27977",
		},
		{
			name:      "success 5.0 rpm container and custom registry",
			registry:  "mysccregistry.com",
			fullImage: "localhost/suse/manager/5.0/x86_64/proxy-helm:latest",
			expected:  "mysccregistry.com/a/150158/27977/suse/manager/5.0/x86_64/proxy-helm:latest-ptf-27977",
		},

		// Success cases - 5.1 Multi-Linux Manager
		{
			name:      "success 5.1 container with 5.1 registry",
			registry:  defaultRegistry51,
			fullImage: defaultRegistry51 + "/proxy-tftpd:5.1.0",
			expected:  "registry.suse.com/a/150158/27977/suse/multi-linux-manager/5.1/x86_64/proxy-tftpd:latest-ptf-27977",
		},

		// Failure cases
		{
			name:                 "fail invalid image",
			registry:             baseRegistryHost,
			fullImage:            "some.domain.com/not/matching/suse/proxy-helm:latest",
			expectedErrorMessage: "invalid image name: some.domain.com/not/matching/suse/proxy-helm:latest",
		},
		{
			name:      "fail 5.0 container and invalid custom registry",
			registry:  "mysccregistry.com/invalid/path",
			fullImage: defaultRegistry50 + "/proxy-helm:latest",
			expectedErrorMessage: "image path 'suse/manager/5.0/x86_64/proxy-helm:' does not start with registry " +
				"path 'invalid/path'",
		},
		{
			name:      "fail 5.0 container with 5.1 registry",
			registry:  defaultRegistry51,
			fullImage: defaultRegistry50 + "/proxy-salt-broker:5.0.0",
			expectedErrorMessage: "image path 'suse/manager/5.0/x86_64/proxy-salt-broker:' does not start with " +
				"registry path 'suse/multi-linux-manager/5.1/x86_64'",
		},
		{
			name:      "fail 5.1 container with 5.0 registry",
			registry:  defaultRegistry50,
			fullImage: defaultRegistry51 + "/proxy-squid:5.0.0",
			expectedErrorMessage: "image path 'suse/multi-linux-manager/5.1/x86_64/proxy-squid:' does not start with" +
				" registry path 'suse/manager/5.0/x86_64'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptfID := defaultPtfID
			if tt.ptfID != "" {
				ptfID = tt.ptfID
			}
			user := defaultUser
			if tt.user != "" {
				user = tt.user
			}
			suffix := defaultSuffix
			if tt.suffix != "" {
				suffix = tt.suffix
			}

			actual, err := ComputePTF(tt.registry, user, ptfID, tt.fullImage, suffix)
			if err != nil {
				if tt.expectedErrorMessage == "" {
					t.Errorf("Unexpected error while executing ComputePTF('%s', '%s', '%s', '%s', '%s'): %s",
						tt.registry, tt.user, tt.ptfID, tt.fullImage, tt.suffix, err)
				} else if !strings.Contains(err.Error(), tt.expectedErrorMessage) {
					t.Errorf("Expected error message to contain '%s', but got: %s",
						tt.expectedErrorMessage, err.Error())
				}
			} else if actual != tt.expected {
				t.Errorf("ComputePTF('%s', '%s', '%s', '%s', '%s') = %s\nexpected: %s",
					tt.registry, tt.user, tt.ptfID, tt.fullImage, tt.suffix, actual, tt.expected)
			}
		})
	}
}

func TestComputeImage(t *testing.T) {
	data := [][]string{
		{"registry:5000/path/to/image:foo", "registry:5000/path/to/image:foo", "bar", ""},
		{"registry:5000/path/to/image:foo", "REGISTRY:5000/path/to/image:foo", "bar", ""},
		{"registry:5000/path/to/image:foo", "REGISTRY:5000/path/to/image:foo", "BAR", ""},
		{"registry:5000/path/to/image:bar", "registry:5000/path/to/image", "bar", ""},
		{"registry/path/to/image:foo", "registry/path/to/image:foo", "bar", ""},
		{"registry/path/to/image:bar", "registry/path/to/image", "bar", ""},
		{"registry/path/to/image:bar", "path/to/image", "bar", "registry"},
		{"registry:5000/path/to/image:foo", "path/to/image:foo", "BAR", "REGISTRY:5000"},
		{"registry:5000/path/to/image-migration-14-16:foo", "registry:5000/path/to/image:foo", "bar", "", "-migration-14-16"},
		{"registry:5000/path/to/image-migration-14-16:bar", "registry:5000/path/to/image", "bar", "", "-migration-14-16"},
		{"registry/path/to/image-migration-14-16:foo", "registry/path/to/image:foo", "bar", "", "-migration-14-16"},
		{"registry/path/to/image-migration-14-16:bar", "registry/path/to/image", "bar", "", "-migration-14-16"},
		{"registry/path/to/image-migration-14-16:bar", "path/to/image", "bar", "registry", "-migration-14-16"},
		{
			// bsc#1226436
			"registry.suse.de/suse/sle-15-sp6/update/products/manager50/containerfile/suse/manager/5.0/x86_64/server:bar",
			"registry.suse.com/suse/manager/5.0/x86_64/server",
			"bar",
			"registry.suse.de/suse/sle-15-sp6/update/products/manager50/containerfile",
			"",
		},
		{
			"cloud.com/suse/manager/5.0/x86_64/server:5.0.0",
			"registry.suse.com/suse/manager/5.0/x86_64/server",
			"5.0.0",
			"cloud.com",
			"",
		},
		{
			"cloud.com/suse/manager/5.0/x86_64/server:5.0.0",
			"/suse/manager/5.0/x86_64/server",
			"5.0.0",
			"cloud.com",
			"",
		},
		{
			"cloud.com/suse/manager/5.0/x86_64/server:5.0.0",
			"suse/manager/5.0/x86_64/server",
			"5.0.0",
			"cloud.com",
			"",
		},
		{
			"cloud.com/my/path/server:5.0.0",
			"my/path/server",
			"5.0.0",
			"cloud.com",
			"",
		},
	}

	for i, testCase := range data {
		result := testCase[0]
		image := types.ImageFlags{
			Name: testCase[1],
			Tag:  testCase[2],
		}
		appendToImage := testCase[4:]

		actual, err := ComputeImage(testCase[3], "defaulttag", image, appendToImage...)

		if err != nil {
			t.Errorf(
				"Testcase %d: Unexpected error while computing image with %s, %s, %s: %s",
				i, image.Name, image.Tag, appendToImage, err,
			)
		}
		if actual != result {
			t.Errorf(
				"Testcase %d: Expected %s got %s when computing image with %s, %s, %s",
				i, result, actual, image.Name, image.Tag, appendToImage,
			)
		}
	}
}

func TestIsWellFormedFQDN(t *testing.T) {
	data := []string{
		"manager.mgr.suse.de",
		"suma50.suse.de",
	}

	for i, testCase := range data {
		if !IsWellFormedFQDN(testCase) {
			t.Errorf("Testcase %d: Unexpected failure while validating FQDN with %s", i, testCase)
		}
	}
	wrongData := []string{
		"manager",
		"suma50",
		"test24.example24.com..",
		"127.0.0.1",
	}

	for i, testCase := range wrongData {
		if IsWellFormedFQDN(testCase) {
			t.Errorf("Testcase %d: Unexpected success while validating FQDN with %s", i, testCase)
		}
	}
}
func TestComputeImageError(t *testing.T) {
	data := [][]string{
		{"registry:path/to/image:tag:tag", "bar"},
	}

	for _, testCase := range data {
		image := types.ImageFlags{
			Name: testCase[0],
			Tag:  testCase[1],
		}

		_, err := ComputeImage("defaultregistry", "defaulttag", image)
		if err == nil {
			t.Errorf("Expected error for %s with tag %s, got none", image.Name, image.Tag)
		}
	}
}

func TestConfig(t *testing.T) {
	type fakeFlags struct {
		firstConf  string
		secondConf string
		thirdConf  string
		fourthConf string
	}
	fakeCmd := &cobra.Command{
		Use:  "podman",
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags fakeFlags
			flags.firstConf = ""
			flags.secondConf = ""
			flags.thirdConf = ""
			flags.fourthConf = ""
			return CommandHelper(nil, cmd, args, &flags, nil, nil)
		},
	}

	fakeCmd.Flags().String("firstConf", "hardcodedDefault", "")
	fakeCmd.Flags().String("secondConf", "hardcodedDefault", "")
	fakeCmd.Flags().String("thirdConf", "hardcodedDefault", "")
	fakeCmd.Flags().String("fourthConf", "hardcodedDefault", "")

	viper, err := ReadConfig(fakeCmd, "conf_test/firstConfFile.yaml", "conf_test/secondConfFile.yaml")
	if err != nil {
		t.Errorf("Unexpected error while reading configuration files: %s", err)
	}

	// This value is not set by conf file, so it should be the hardcoded default value
	if viper.Get("firstConf") != "hardcodedDefault" {
		t.Errorf("firstConf is %s, instead of hardcodedDefault", viper.Get("firstConf"))
	}
	// This value is set by firstConfFile.yaml
	if viper.Get("secondConf") != "firstConfFile" {
		t.Errorf("secondConf is %s, instead of firstConfFile", viper.Get("secondConf"))
	}
	// This value is as first set by firstConfFile.yaml, but then overwritten by secondConfFile.yaml
	if viper.Get("thirdConf") != "SecondConfFile" {
		t.Errorf("thirdConf is %s, instead of SecondConfFile", viper.Get("thirdConf"))
	}
	// This value is set by secondConfFile.yaml
	if viper.Get("fourthconf") != "SecondConfFile" {
		t.Errorf("fourthconf is %s, instead of SecondConfFile", viper.Get("fourthconf"))
	}
}

// Test saveBinaryData function.
func TestSaveBinaryData(t *testing.T) {
	testDir := t.TempDir()

	filepath := path.Join(testDir, "testfile")
	data := []int8{104, 101, 121, 32, 116, 104, 101, 114, 101, 33}

	// Save binary data to a file
	err := SaveBinaryData(filepath, data)
	testutils.AssertTrue(t, "Unexpected error executing SaveBinaryData", err == nil)

	// Read the file back and compare contents
	storedData := testutils.ReadFileAsBinary(t, filepath)
	testutils.AssertEquals(
		t, "File configuration binary doesn't match",
		fmt.Sprintf("%v", data), fmt.Sprintf("%v", storedData),
	)
}

func TestCompareVersion(t *testing.T) {
	testutils.AssertTrue(t, "2024.07 is not inferior to 2024.13", CompareVersion("2024.07", "2024.13") < 0)
	testutils.AssertTrue(t, "2024.13 is not superior to 2024.07", CompareVersion("2024.13", "2024.07") > 0)
	testutils.AssertTrue(t, "2024.13 is not equal to 2024.13", CompareVersion("2024.13", "2024.13") == 0)
}

// TestRemoveRegistryFromImage tests the RemoveRegistryFromImage function.
func TestRemoveRegistryFromImage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "regular registry",
			input:    "registry.example.com/some/name/space/image",
			expected: "some/name/space/image",
		},
		{
			name:     "url with protocol",
			input:    "https://registry.example.com/some/name/space/image",
			expected: "some/name/space/image",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only registry fqdn",
			input:    "registry.example.com",
			expected: "",
		},
		{
			name:     "registry with port",
			input:    "registry:5000/no/problemo/namespace/image",
			expected: "no/problemo/namespace/image",
		},
		{
			name:     "path only",
			input:    "path/only/image",
			expected: "path/only/image",
		},
		{
			name:     "single component",
			input:    "image",
			expected: "image",
		},
		{
			name:     "protocol with path",
			input:    "http://registry.example.com/path/to/image",
			expected: "path/to/image",
		},
		{
			name:     "complex registry with multiple dots",
			input:    "my.complex.registry.com/org/project/image",
			expected: "org/project/image",
		},
		{
			name:     "localhost registry",
			input:    "localhost:5000/my/image",
			expected: "my/image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := RemoveRegistryFromImage(tt.input)

			if actual != tt.expected {
				t.Errorf("RemoveRegistryFromImage(%q) = %q, expected %q", tt.input, actual, tt.expected)
			}
		})
	}
}

func TestSplitRegistryHostAndPath(t *testing.T) {
	// Constants
	const (
		exampleRegistry      = "registry.example.com"
		localhostWithPort    = "localhost:5000"
		suseRegistryWithPort = "registry.suse.com:443"
		complexRegistry      = "my.complex.registry.com"
	)

	tests := []struct {
		name         string
		registry     string
		expectedHost string
		expectedPath string
	}{
		{
			name:         "registry with path",
			registry:     exampleRegistry + "/path/to/namespace",
			expectedHost: exampleRegistry,
			expectedPath: "path/to/namespace",
		},
		{
			name:         "registry with single path component",
			registry:     exampleRegistry + "/namespace",
			expectedHost: exampleRegistry,
			expectedPath: "namespace",
		},
		{
			name:         "registry without path",
			registry:     exampleRegistry,
			expectedHost: exampleRegistry,
			expectedPath: "",
		},
		{
			name:         "localhost registry with port and path",
			registry:     localhostWithPort + "/my/namespace",
			expectedHost: localhostWithPort,
			expectedPath: "my/namespace",
		},
		{
			name:         "localhost registry with port only",
			registry:     localhostWithPort,
			expectedHost: localhostWithPort,
			expectedPath: "",
		},
		{
			name:         "registry with port and path",
			registry:     suseRegistryWithPort + "/suse/manager/5.0/x86_64",
			expectedHost: suseRegistryWithPort,
			expectedPath: "suse/manager/5.0/x86_64",
		},
		{
			name:         "registry with port only",
			registry:     suseRegistryWithPort,
			expectedHost: suseRegistryWithPort,
			expectedPath: "",
		},
		{
			name:         "complex path",
			registry:     complexRegistry + "/org/project/sub/namespace",
			expectedHost: complexRegistry,
			expectedPath: "org/project/sub/namespace",
		},
		{
			name:         "empty registry",
			registry:     "",
			expectedHost: "",
			expectedPath: "",
		},
		{
			name:         "registry starting with slash",
			registry:     "/path/only",
			expectedHost: "",
			expectedPath: "path/only",
		},
		{
			name:         "registry with multiple slashes",
			registry:     exampleRegistry + "/path/with/many/slashes",
			expectedHost: exampleRegistry,
			expectedPath: "path/with/many/slashes",
		},
		{
			name:         "single component without slash",
			registry:     "localhost",
			expectedHost: "localhost",
			expectedPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualHost, actualPath := SplitRegistryHostAndPath(tt.registry)

			if actualHost != tt.expectedHost {
				t.Errorf("SplitRegistryHostAndPath(%q) host = %q, expected %q", tt.registry, actualHost, tt.expectedHost)
			}
			if actualPath != tt.expectedPath {
				t.Errorf("SplitRegistryHostAndPath(%q) path = %q, expected %q", tt.registry, actualPath, tt.expectedPath)
			}
		})
	}
}
