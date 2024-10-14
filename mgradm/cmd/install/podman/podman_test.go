// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils/flags_tests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := flags_tests.InstallFlagsTestArgs()
	args = append(args, flags_tests.PodmanFlagsTestArgs...)
	args = append(args, "srv.fq.dn")

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *podmanInstallFlags,
		cmd *cobra.Command, args []string,
	) error {
		flags_tests.AssertInstallFlags(t, cmd, &flags.InstallFlags)
		flags_tests.AssertPodmanInstallFlags(t, cmd, &flags.Podman)
		test_utils.AssertEquals(t, "Wrong FQDN", "srv.fq.dn", args[0])
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

func TestParamsChangedConfig(t *testing.T) {
	config := `
coco:
  replicas: 2
hubxmlrpc:
  replicas: 0`

	dir := t.TempDir()
	configPath := path.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(config), 0600); err != nil {
		t.Fatalf("Failed to write config file: %s", err)
	}

	tester := func(globalFlags *types.GlobalFlags, flags *podmanInstallFlags,
		cmd *cobra.Command, args []string,
	) error {
		test_utils.AssertEquals(t, "Coco replicas badly parsed", 2, flags.Coco.Replicas)
		test_utils.AssertTrue(t, "Coco replicas not marked as changed", flags.Coco.IsChanged)
		test_utils.AssertEquals(t, "Hub XML-RPC API replicas badly parsed", 0, flags.HubXmlrpc.Replicas)
		test_utils.AssertTrue(t, "Hub XML-RPC API replicas not marked as changed", flags.HubXmlrpc.IsChanged)
		return nil
	}

	globalFlags := types.GlobalFlags{ConfigPath: configPath}
	cmd := newCmd(&globalFlags, tester)

	cmd.SetArgs([]string{"srv.fq.dn"})
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}

func TestParamsNoConfig(t *testing.T) {
	tester := func(globalFlags *types.GlobalFlags, flags *podmanInstallFlags,
		cmd *cobra.Command, args []string,
	) error {
		test_utils.AssertEquals(t, "Coco replicas badly parsed", 0, flags.Coco.Replicas)
		test_utils.AssertTrue(t, "Coco replicas marked as changed", !flags.Coco.IsChanged)
		test_utils.AssertEquals(t, "Hub XML-RPC API replicas badly parsed", 0, flags.HubXmlrpc.Replicas)
		test_utils.AssertTrue(t, "Hub XML-RPC API replicas marked as changed", !flags.HubXmlrpc.IsChanged)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newCmd(&globalFlags, tester)

	cmd.SetArgs([]string{"srv.fq.dn"})
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}
