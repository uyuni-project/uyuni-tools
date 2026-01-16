// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package upgrade

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := flagstests.ServerFlagsTestArgs()
	args = append(args, flagstests.PodmanFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *podmanUpgradeFlags,
		_ *cobra.Command, _ []string,
	) error {
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

func TestListParamsParsing(t *testing.T) {
	args := []string{}

	args = append(args, flagstests.ImageFlagsTestArgs...)
	args = append(args, flagstests.SCCFlagTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(flags *podmanUpgradeFlags) error {
		flagstests.AssertImageFlag(t, &flags.Image)
		flagstests.AssertSCCFlag(t, &flags.Installation.SCC)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newListCmd(&globalFlags, tester)

	testutils.AssertHasAllFlags(t, cmd, args)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}
