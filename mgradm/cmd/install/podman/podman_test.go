// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := flagstests.InstallFlagsTestArgs()
	args = append(args, flagstests.MirrorFlagTestArgs...)
	args = append(args, flagstests.PodmanFlagsTestArgs...)
	args = append(args, "srv.fq.dn")

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *podmanInstallFlags,
		_ *cobra.Command, args []string,
	) error {
		flagstests.AssertMirrorFlag(t, flags.Mirror)
		flagstests.AssertInstallFlags(t, &flags.ServerFlags)
		flagstests.AssertPodmanInstallFlags(t, &flags.Podman)
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

func TestParamsChangedConfig(t *testing.T) {
	config := `
coco:
  replicas: 2
hubxmlrpc:
  replicas: 0
saline:
  port: 8226
  replicas: 1`

	dir := t.TempDir()
	configPath := path.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(config), 0600); err != nil {
		t.Fatalf("Failed to write config file: %s", err)
	}

	tester := func(_ *types.GlobalFlags, flags *podmanInstallFlags,
		_ *cobra.Command, _ []string,
	) error {
		testutils.AssertEquals(t, "Coco replicas badly parsed", 2, flags.Coco.Replicas)
		testutils.AssertTrue(t, "Coco replicas not marked as changed", flags.Coco.IsChanged)
		testutils.AssertEquals(t, "Hub XML-RPC API replicas badly parsed", 0, flags.HubXmlrpc.Replicas)
		testutils.AssertTrue(t, "Hub XML-RPC API replicas not marked as changed", flags.HubXmlrpc.IsChanged)
		testutils.AssertEquals(t, "Saline replicas badly parsed", 1, flags.Saline.Replicas)
		testutils.AssertEquals(t, "Saline port badly parsed", 8226, flags.Saline.Port)
		testutils.AssertTrue(t, "Saline flags not marked as changed", flags.Saline.IsChanged)
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
	tester := func(_ *types.GlobalFlags, flags *podmanInstallFlags,
		_ *cobra.Command, _ []string,
	) error {
		testutils.AssertEquals(t, "Coco replicas badly parsed", 0, flags.Coco.Replicas)
		testutils.AssertTrue(t, "Coco replicas marked as changed", !flags.Coco.IsChanged)
		testutils.AssertEquals(t, "Hub XML-RPC API replicas badly parsed", 0, flags.HubXmlrpc.Replicas)
		testutils.AssertTrue(t, "Hub XML-RPC API replicas marked as changed", !flags.HubXmlrpc.IsChanged)
		testutils.AssertEquals(t, "Saline replicas badly parsed", 0, flags.Saline.Replicas)
		testutils.AssertEquals(t, "Saline port badly parsed", 8216, flags.Saline.Port)
		testutils.AssertTrue(t, "Saline flags marked as changed", !flags.Saline.IsChanged)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newCmd(&globalFlags, tester)

	cmd.SetArgs([]string{"srv.fq.dn"})
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}
