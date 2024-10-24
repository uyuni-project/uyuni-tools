// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
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
		testutils.AssertEquals(t, "Error parsing --dababase", "reportdb", flags.Database)
		testutils.AssertTrue(t, "Error parsing --interactive", flags.Interactive)
		testutils.AssertTrue(t, "Error parsing --force", flags.ForceOverwrite)
		testutils.AssertEquals(t, "Error parsing --dababase", "reportdb", flags.Database)
		testutils.AssertEquals(t, "Error parsing --output", "path/to/output", flags.OutputFile)
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
