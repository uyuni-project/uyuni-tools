// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package migrate

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--prepare",
		"--user", "sudoer",
		"source.fq.dn",
	}

	args = append(args, flagstests.MirrorFlagTestArgs...)
	args = append(args, flagstests.SCCFlagTestArgs...)
	args = append(args, flagstests.PodmanFlagsTestArgs...)
	args = append(args, flagstests.ServerFlagsTestArgs()...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *podmanMigrateFlags,
		_ *cobra.Command, args []string,
	) error {
		testutils.AssertTrue(t, "Prepare not set", flags.Migration.Prepare)
		flagstests.AssertMirrorFlag(t, flags.Mirror)
		testutils.AssertEquals(t, "Error parsing --user", "sudoer", flags.Migration.User)
		testutils.AssertEquals(t, "Wrong FQDN", "source.fq.dn", args[0])
		flagstests.AssertPodmanInstallFlags(t, &flags.Podman)
		flagstests.AssertServerFlags(t, &flags.ServerFlags)
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
