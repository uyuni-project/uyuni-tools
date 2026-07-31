// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package addca

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{}
	args = append(args, flagstests.SSLGenerationFlagsTestArgs...)
	args = append(args,
		"--ssl-ca-root", "path/root.crt",
		"--ssl-db-ca-root", "path/db-root.crt",
		"--ssl-password", "sslsecret",
	)

	// Test function asserting that the args are properly parsed.
	tester := func(_ *types.GlobalFlags, flags *addCAFlags, _ *cobra.Command, _ []string) error {
		flagstests.AssertSSLGenerationFlag(t, &flags.SSL.SSLCertGenerationFlags)
		testutils.AssertEquals(t, "Error parsing --ssl-ca-root", "path/root.crt", flags.SSL.Ca.Root)
		testutils.AssertEquals(t, "Error parsing --ssl-db-ca-root", "path/db-root.crt", flags.SSL.DB.CA.Root)
		testutils.AssertEquals(t, "Error parsing --ssl-password", "sslsecret", flags.SSL.Password)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newCmd(&globalFlags, tester)

	testutils.AssertHasAllFlags(t, cmd, args)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}
