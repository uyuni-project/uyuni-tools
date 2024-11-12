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
	args := []string{}

	args = append(args, flagstests.ImageFlagsTestArgs...)
	args = append(args, flagstests.DBUpdateImageFlagTestArgs...)
	args = append(args, flagstests.CocoFlagsTestArgs...)
	args = append(args, flagstests.HubXmlrpcFlagsTestArgs...)
	args = append(args, flagstests.SCCFlagTestArgs...)
	args = append(args, flagstests.PodmanFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *podmanUpgradeFlags,
		cmd *cobra.Command, args []string,
	) error {
		flagstests.AssertImageFlag(t, cmd, &flags.Image)
		flagstests.AssertDBUpgradeImageFlag(t, cmd, &flags.DBUpgradeImage)
		flagstests.AssertCocoFlag(t, cmd, &flags.Coco)
		flagstests.AssertHubXmlrpcFlag(t, cmd, &flags.HubXmlrpc)
		flagstests.AssertSCCFlag(t, cmd, &flags.SCC)
		flagstests.AssertPodmanInstallFlags(t, cmd, &flags.Podman)
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
