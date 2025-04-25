// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package gpglist

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--system",
		"--backend", "kubectl",
	}

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *gpgListFlags, _ *cobra.Command, _ []string) error {
		testutils.AssertTrue(t, "Error parsing --system", flags.System)
		testutils.AssertEquals(t, "Error parsing --backend", "kubectl", flags.Backend)
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
