// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := flagstests.MigrateFlagsTestArgs()
	args = append(args, flagstests.PodmanFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *podmanMigrateFlags,
		_ *cobra.Command, args []string,
	) error {
		testutils.AssertEquals(t, "Wrong FQDN", "source.fq.dn", args[0])
		flagstests.AssertServerFlags(t, &flags.ServerFlags)
		flagstests.AssertMigrateFlags(t, &flags.MigrationFlags)
		flagstests.AssertUpgradeFlags(t, &flags.UpgradeFlags)
		flagstests.AssertPodmanInstallFlags(t, &flags.Podman)
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
