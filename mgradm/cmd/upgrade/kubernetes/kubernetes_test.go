// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils/flags_tests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{}

	args = append(args, flags_tests.ImageFlagsTestArgs...)
	args = append(args, flags_tests.DbUpdateImageFlagTestArgs...)
	args = append(args, flags_tests.CocoFlagsTestArgs...)
	args = append(args, flags_tests.HubXmlrpcFlagsTestArgs...)
	args = append(args, flags_tests.SccFlagTestArgs...)
	args = append(args, flags_tests.ServerHelmFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *kubernetesUpgradeFlags,
		cmd *cobra.Command, args []string,
	) error {
		flags_tests.AssertImageFlag(t, cmd, &flags.Image)
		flags_tests.AssertDbUpgradeImageFlag(t, cmd, &flags.DbUpgradeImage)
		flags_tests.AssertCocoFlag(t, cmd, &flags.Coco)
		flags_tests.AssertHubXmlrpcFlag(t, cmd, &flags.HubXmlrpc)
		// TODO Assert SCC flags
		flags_tests.AssertServerHelmFlags(t, cmd, &flags.Helm)
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
