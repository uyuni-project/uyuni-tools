// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package register

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils/flags_tests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--backend", "kubectl",
	}

	args = append(args, flags_tests.APIFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *registerFlags,
		cmd *cobra.Command, args []string,
	) error {
		test_utils.AssertEquals(t, "Error parsing --backend", "kubectl", flags.Backend)
		flags_tests.AssertAPIFlags(t, cmd, &flags.ConnectionDetails)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newCmd(&globalFlags, tester)

	test_utils.AssertHasAllFlags(t, cmd, args)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}
