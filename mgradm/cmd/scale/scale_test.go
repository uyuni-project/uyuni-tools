// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package scale

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--replicas", "2",
		"some-service",
	}

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *scaleFlags, _ *cobra.Command, args []string) error {
		testutils.AssertEquals(t, "Error parsing --replicas", 2, flags.Replicas)
		testutils.AssertEquals(t, "Error parsing the service name", "some-service", args[0])
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
