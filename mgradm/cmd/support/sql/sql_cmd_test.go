// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--database", "reportdb",
		"--interactive",
		"--force",
		"--output", "path/to/output",
		"--backend", "kubectl",
	}

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *sqlFlags,
		cmd *cobra.Command, args []string,
	) error {
		test_utils.AssertEquals(t, "Error parsing --dababase", "reportdb", flags.Database)
		test_utils.AssertTrue(t, "Error parsing --interactive", flags.Interactive)
		test_utils.AssertTrue(t, "Error parsing --force", flags.ForceOverwrite)
		test_utils.AssertEquals(t, "Error parsing --dababase", "reportdb", flags.Database)
		test_utils.AssertEquals(t, "Error parsing --output", "path/to/output", flags.OutputFile)
		test_utils.AssertEquals(t, "Error parsing --backend", "kubectl", flags.Backend)
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
