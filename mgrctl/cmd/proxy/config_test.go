// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func getCommonArgs() []string {
	args := []string{
		"--proxy-name", "pxy1.test.com",
		"--proxy-sshPort", "1234",
		"--proxy-parent", "uyuni.test.com",
		"--proxy-maxCache", "123456",
		"--proxy-email", "admin@proxy.test.com",
		"--output", "path/to/output.tgz",
		"--ssl-ca-cert", "path/to/ca.crt",
	}
	args = append(args, flagstests.APIFlagsTestArgs...)
	return args
}

func assertCommonArgs(t *testing.T, cmd *cobra.Command, flags *proxyCreateConfigFlags) {
	flagstests.AssertAPIFlags(t, cmd, &flags.ConnectionDetails)
	testutils.AssertEquals(t, "Unexpected proxy name", "pxy1.test.com", flags.Proxy.Name)
	testutils.AssertEquals(t, "Unexpected proxy SSH port", 1234, flags.Proxy.Port)
	testutils.AssertEquals(t, "Unexpected proxy parent", "uyuni.test.com", flags.Proxy.Parent)
	testutils.AssertEquals(t, "Unexpected proxy max cache", 123456, flags.Proxy.MaxCache)
	testutils.AssertEquals(t, "Unexpected proxy email", "admin@proxy.test.com", flags.Proxy.Email)
	testutils.AssertEquals(t, "Unexpected output path", "path/to/output.tgz", flags.Output)
	testutils.AssertEquals(t, "Unexpected SSL CA cert path", "path/to/ca.crt", flags.SSL.Ca.Cert)
}

func TestParamsParsingGeneratedCert(t *testing.T) {
	args := getCommonArgs()
	args = append(args,
		"--ssl-ca-cert", "path/to/ca.crt",
		"--ssl-ca-key", "path/to/ca.key",
		"--ssl-ca-password", "casecret",
		"--ssl-email", "ssl@test.com",
	)
	args = append(args, flagstests.SSLGenerationFlagsTestArgs...)

	conflictingFlags := []string{
		"--ssl-proxy-cert",
		"--ssl-proxy-key",
		"--ssl-ca-intermediate",
	}

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *proxyCreateConfigFlags,
		cmd *cobra.Command, args []string,
	) error {
		assertCommonArgs(t, cmd, flags)
		flagstests.AssertSSLGenerationFlags(t, cmd, &flags.SSL.SSLCertGenerationFlags)
		testutils.AssertEquals(t, "Unexpected SSL CA cert path", "path/to/ca.crt", flags.SSL.Ca.Cert)
		testutils.AssertEquals(t, "Unexpected SSL CA key path", "path/to/ca.key", flags.SSL.Ca.Key)
		testutils.AssertEquals(t, "Unexpected SSL CA password", "casecret", flags.SSL.Ca.Password)
		testutils.AssertEquals(t, "Unexpected SSL email", "ssl@test.com", flags.SSL.Email)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newCmd(&globalFlags, tester)

	testutils.AssertHasAllFlagsIgnores(t, cmd, args, conflictingFlags)

	t.Logf("flags: %s", strings.Join(args, " "))
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}

func TestParamsParsingProvidedCert(t *testing.T) {
	args := getCommonArgs()
	args = append(args,
		"--ssl-ca-intermediate", "path/to/ca1.crt",
		"--ssl-ca-intermediate", "path/to/ca2.crt",
		"--ssl-proxy-cert", "path/to/proxy.crt",
		"--ssl-proxy-key", "path/to/proxy.key",
	)

	conflictingFlags := []string{
		"--ssl-email",
		"--ssl-ca-key",
		"--ssl-ca-password",
		"--ssl-cname",
		"--ssl-country",
		"--ssl-state",
		"--ssl-city",
		"--ssl-org",
		"--ssl-ou",
	}

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *proxyCreateConfigFlags,
		cmd *cobra.Command, args []string,
	) error {
		assertCommonArgs(t, cmd, flags)
		testutils.AssertEquals(t, "Unexpected SSL CA cert path", "path/to/ca.crt", flags.SSL.Ca.Cert)
		testutils.AssertEquals(t, "Unexpected SSL intermediate CA cert paths",
			[]string{"path/to/ca1.crt", "path/to/ca2.crt"}, flags.SSL.Ca.Intermediate,
		)
		testutils.AssertEquals(t, "Unexpected Proxy SSL cert path", "path/to/proxy.crt", flags.SSL.Proxy.Cert)
		testutils.AssertEquals(t, "Unexpected Proxy SSL key path", "path/to/proxy.key", flags.SSL.Proxy.Key)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newCmd(&globalFlags, tester)

	testutils.AssertHasAllFlagsIgnores(t, cmd, args, conflictingFlags)

	t.Logf("flags: %s", strings.Join(args, " "))
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}
