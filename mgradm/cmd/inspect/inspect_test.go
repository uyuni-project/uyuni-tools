// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{}

	args = append(args, flagstests.ImageFlagsTestArgs...)
	args = append(args, flagstests.SCCFlagTestArgs...)
	args = append(args, flagstests.PgsqlFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *inspectFlags, _ *cobra.Command, _ []string) error {
		flagstests.AssertImageFlag(t, &flags.Image)
		flagstests.AssertRegistryFlag(t, &flags.Image.Registry)
		flagstests.AssertSCCFlag(t, &flags.SCC)
		flagstests.AssertPgsqlFlag(t, &flags.Pgsql)
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
