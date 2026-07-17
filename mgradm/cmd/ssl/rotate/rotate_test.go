// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package rotate

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
	args = append(args, flagstests.InstallSSLFlagsTestArgs...)
	args = append(args, flagstests.InstallDBSSLFlagsTestArgs...)
	args = append(args, "--force", "--check-only", "--emergency")

	// Test function asserting that the args are properly parsed.
	tester := func(_ *types.GlobalFlags, flags *rotateFlags, cmd *cobra.Command, _ []string) error {
		flagstests.AssertSSLGenerationFlag(t, &flags.SSL.SSLCertGenerationFlags)
		flagstests.AssertInstallSSLFlag(t, &flags.SSL)
		flagstests.AssertInstallDBSSLFlag(t, &flags.SSL.DB)
		testutils.AssertTrue(t, "Error parsing --force", flags.Force)
		testutils.AssertTrue(t, "Error parsing --emergency", flags.Emergency)
		checkOnly, _ := cmd.Flags().GetBool("check-only")
		testutils.AssertTrue(t, "Error parsing --check-only", checkOnly)
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
