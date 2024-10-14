// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package scale

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--replicas", "2",
		"--backend", "kubectl",
		"some-service",
	}

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *scaleFlags,
		cmd *cobra.Command, args []string,
	) error {
		test_utils.AssertEquals(t, "Error parsing --replicas", 2, flags.Replicas)
		test_utils.AssertEquals(t, "Error parsing --backend", "kubectl", flags.Backend)
		test_utils.AssertEquals(t, "Error parsing the service name", "some-service", args[0])
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
