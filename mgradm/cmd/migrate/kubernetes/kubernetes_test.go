// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := flagstests.MigrateFlagsTestArgs()

	args = append(args, flagstests.ServerKubernetesFlagsTestArgs...)
	args = append(args, flagstests.VolumesFlagsTestExpected...)
	args = append(args, "--ssh-key-public", "path/ssh.pub")
	args = append(args, "--ssh-key-private", "path/ssh")
	args = append(args, "--ssh-key-public", "path/ssh.pub")
	args = append(args, "--ssh-knownhosts", "path/known_hosts")
	args = append(args, "--ssh-config", "path/config")

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *kubernetes.KubernetesServerFlags,
		_ *cobra.Command, args []string,
	) error {
		testutils.AssertEquals(t, "Wrong FQDN", "source.fq.dn", args[0])
		flagstests.AssertMirrorFlag(t, flags.Mirror)
		flagstests.AssertMigrateFlags(t, &flags.Migration)
		flagstests.AssertUpgradeFlags(t, &flags.Upgrade)
		flagstests.AssertVolumesFlags(t, &flags.Volumes)

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
