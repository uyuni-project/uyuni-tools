// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--force",
		"--purge-volumes",
		"--purge-images",
		"--backend", "kubectl",
	}

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *utils.UninstallFlags,
		cmd *cobra.Command, args []string,
	) error {
		testutils.AssertTrue(t, "Error parsing --force", flags.Force)
		testutils.AssertTrue(t, "Error parsing --purge-volumes", flags.Purge.Volumes)
		testutils.AssertTrue(t, "Error parsing --purge-images", flags.Purge.Images)
		testutils.AssertEquals(t, "Error parsing --backend", "kubectl", flags.Backend)
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
