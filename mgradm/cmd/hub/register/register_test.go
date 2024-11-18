// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package register

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestParamsParsing(t *testing.T) {
	args := []string{}
	if utils.KubernetesBuilt {
		args = append(args, "--backend", "kubectl")
	}

	args = append(args, flagstests.APIFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *registerFlags, _ *cobra.Command, _ []string) error {
		if utils.KubernetesBuilt {
			testutils.AssertEquals(t, "Error parsing --backend", "kubectl", flags.Backend)
		}
		flagstests.AssertAPIFlags(t, &flags.ConnectionDetails)
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
