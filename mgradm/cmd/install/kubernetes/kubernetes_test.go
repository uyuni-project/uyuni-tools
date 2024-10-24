// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flags_tests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := flags_tests.InstallFlagsTestArgs()
	args = append(args, flags_tests.ServerHelmFlagsTestArgs...)
	args = append(args, "srv.fq.dn")

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *kubernetesInstallFlags,
		cmd *cobra.Command, args []string,
	) error {
		flags_tests.AssertInstallFlags(t, cmd, &flags.InstallFlags)
		flags_tests.AssertServerHelmFlags(t, cmd, &flags.Helm)
		testutils.AssertEquals(t, "Wrong FQDN", "srv.fq.dn", args[0])
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
