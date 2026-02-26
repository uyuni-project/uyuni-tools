// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package ptf

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--ptf", "ptf123",
		"--test", "test123",
		"--user", "sccuser",
	}

	args = append(args, flagstests.SCCFlagTestArgs...)
	args = append(args, flagstests.ImageProxyFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *podmanPTFFlags,
		_ *cobra.Command, _ []string,
	) error {
		flagstests.AssertSCCFlag(t, &flags.UpgradeFlags.SCC)
		flagstests.AssertProxyImageFlags(t, &flags.UpgradeFlags.ProxyImageFlags)
		testutils.AssertEquals(t, "Error parsing --ptf", "ptf123", flags.PTFId)
		testutils.AssertEquals(t, "Error parsing --test", "test123", flags.TestID)
		testutils.AssertEquals(t, "Error parsing --user", "sccuser", flags.CustomerID)
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
