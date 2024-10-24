// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flags_tests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--backend", "kubectl",
	}

	args = append(args, flags_tests.ImageFlagsTestArgs...)
	args = append(args, flags_tests.SccFlagTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *inspectFlags,
		cmd *cobra.Command, args []string,
	) error {
		flags_tests.AssertImageFlag(t, cmd, &flags.Image)
		flags_tests.AssertSccFlag(t, cmd, &flags.SCC)
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
