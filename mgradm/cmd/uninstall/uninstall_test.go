// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--force",
		"--purge-volumes",
		"--purge-images",
	}
	if utils.KubernetesBuilt {
		args = append(args, "--backend", "kubectl")
	}

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *utils.UninstallFlags,
		cmd *cobra.Command, args []string,
	) error {
		test_utils.AssertTrue(t, "Error parsing --force", flags.Force)
		test_utils.AssertTrue(t, "Error parsing --purge-volumes", flags.Purge.Volumes)
		test_utils.AssertTrue(t, "Error parsing --purge-images", flags.Purge.Images)
		if utils.KubernetesBuilt {
			test_utils.AssertEquals(t, "Error parsing --backend", "kubectl", flags.Backend)
		}
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newCmd(&globalFlags, tester)

	test_utils.AssertHasAllFlags(t, cmd, args)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}
