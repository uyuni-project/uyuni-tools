// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package podman

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{}

	args = append(args, flagstests.SCCFlagTestArgs...)
	args = append(args, flagstests.ImageProxyFlagsTestArgs...)
	args = append(args, flagstests.PodmanFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *podman.PodmanProxyFlags,
		cmd *cobra.Command, args []string,
	) error {
		flagstests.AssertSCCFlag(t, cmd, &flags.SCC)
		flagstests.AssertPodmanInstallFlags(t, cmd, &flags.Podman)
		flagstests.AssertProxyImageFlags(t, cmd, &flags.ProxyImageFlags)
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
