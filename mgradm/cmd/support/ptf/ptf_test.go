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
		"--pullPolicy", "never",
		"--registry", "myOldRegistry",
		"--registry-host", "myoverwrittenregistry",
		"--registry-user", "user",
		"--registry-password", "password",
	}
	args = append(args, flagstests.SCCFlagTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *podmanPTFFlags, _ *cobra.Command, _ []string) error {
		testutils.AssertEquals(t, "Error parsing --ptf", "ptf123", flags.PTFId)
		testutils.AssertEquals(t, "Error parsing --test", "test123", flags.TestID)
		testutils.AssertEquals(t, "Error parsing --user", "sccuser", flags.CustomerID)
		testutils.AssertEquals(t, "Error parsing --pullPolicy", "never", flags.Image.PullPolicy)
		testutils.AssertEquals(t, "Error parsing --registry", "myOldRegistry", flags.Image.Registry.Host)
		flagstests.AssertRegistryFlag(t, &flags.Image.Registry)
		flagstests.AssertSCCFlag(t, &flags.Installation.SCC)
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
