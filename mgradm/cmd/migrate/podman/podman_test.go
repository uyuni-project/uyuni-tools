// SPDX-FileCopyrightText: 2024 SUSE LLC
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
	args := []string{
		"--prepare",
		"--user", "sudoer",
		"source.fq.dn",
	}

	args = append(args, flagstests.MirrorFlagTestArgs...)
	args = append(args, flagstests.SCCFlagTestArgs...)
	args = append(args, flagstests.ImageFlagsTestArgs...)
	args = append(args, flagstests.DBUpdateImageFlagTestArgs...)
	args = append(args, flagstests.CocoFlagsTestArgs...)
	args = append(args, flagstests.HubXmlrpcFlagsTestArgs...)
	args = append(args, flagstests.SalineFlagsTestArgs...)
	args = append(args, flagstests.PodmanFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *podmanMigrateFlags, _ *cobra.Command, args []string) error {
		testutils.AssertTrue(t, "Prepare not set", flags.Prepare)
		flagstests.AssertMirrorFlag(t, flags.Mirror)
		flagstests.AssertSCCFlag(t, &flags.SCC)
		flagstests.AssertImageFlag(t, &flags.Image)
		flagstests.AssertDBUpgradeImageFlag(t, &flags.DBUpgradeImage)
		flagstests.AssertCocoFlag(t, &flags.Coco)
		flagstests.AssertHubXmlrpcFlag(t, &flags.HubXmlrpc)
		flagstests.AssertSalineFlag(t, &flags.Saline)
		testutils.AssertEquals(t, "Error parsing --user", "sudoer", flags.User)
		flagstests.AssertPodmanInstallFlags(t, &flags.Podman)
		testutils.AssertEquals(t, "Wrong FQDN", "source.fq.dn", args[0])
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
