// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package logs

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--follow",
		"--timestamps",
		"--tail=20",
		"--since", "3h",
		"container1", "container2",
	}

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *logsFlags,
		_ *cobra.Command, _ []string,
	) error {
		testutils.AssertTrue(t, "Error parsing --follow", flags.Follow)
		testutils.AssertTrue(t, "Error parsing --timestamps", flags.Timestamps)
		testutils.AssertEquals(t, "Error parsing --tail", 20, flags.Tail)
		testutils.AssertEquals(t, "Error parsing --since", "3h", flags.Since)
		testutils.AssertEquals(t, "Error parsing containers", []string{"container1", "container2"}, flags.Containers)
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
