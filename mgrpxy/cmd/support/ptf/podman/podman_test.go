// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package podman

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flags_tests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--ptf", "ptf123",
		"--test", "test123",
		"--user", "sccuser",
	}

	args = append(args, flags_tests.SccFlagTestArgs...)
	args = append(args, flags_tests.ImageProxyFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *podmanPTFFlags,
		cmd *cobra.Command, args []string,
	) error {
		flags_tests.AssertSccFlag(t, cmd, &flags.UpgradeFlags.SCC)
		flags_tests.AssertProxyImageFlags(t, cmd, &flags.UpgradeFlags.ProxyImageFlags)
		testutils.AssertEquals(t, "Error parsing --ptf", "ptf123", flags.PTFId)
		testutils.AssertEquals(t, "Error parsing --test", "test123", flags.TestId)
		testutils.AssertEquals(t, "Error parsing --user", "sccuser", flags.CustomerId)
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
