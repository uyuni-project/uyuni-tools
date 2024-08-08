// SPDX-FileCopyrightText: 2024 SUSE LLC
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
	args := []string{
		"--prepare",
		"--user", "sudoer",
		"--ssl-password", "sslsecret",
		"--ssh-key-public", "path/ssh.pub",
		"--ssh-key-private", "path/ssh",
		"--ssh-knownhosts", "path/known_hosts",
		"--ssh-config", "path/config",
		"source.fq.dn",
	}

	args = append(args, flagstests.MirrorFlagTestArgs...)
	args = append(args, flagstests.SCCFlagTestArgs...)
	args = append(args, flagstests.ImageFlagsTestArgs...)
	args = append(args, flagstests.DBUpdateImageFlagTestArgs...)
	args = append(args, flagstests.CocoFlagsTestArgs...)
	args = append(args, flagstests.HubXmlrpcFlagsTestArgs...)
	args = append(args, flagstests.SalineFlagsTestArgs...)
	args = append(args, flagstests.ServerHelmFlagsTestArgs...)
	args = append(args, flagstests.VolumesFlagsTestExpected...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *kubernetes.KubernetesServerFlags,
		_ *cobra.Command, args []string,
	) error {
		testutils.AssertTrue(t, "Prepare not set", flags.Migration.Prepare)
		flagstests.AssertMirrorFlag(t, flags.Mirror)
		flagstests.AssertSCCFlag(t, &flags.Installation.SCC)
		flagstests.AssertImageFlag(t, &flags.Image)
		flagstests.AssertDBUpgradeImageFlag(t, &flags.DBUpgradeImage)
		flagstests.AssertCocoFlag(t, &flags.Coco)
		flagstests.AssertHubXmlrpcFlag(t, &flags.HubXmlrpc)
		flagstests.AssertSalineFlag(t, &flags.Saline)
		testutils.AssertEquals(t, "Error parsing --user", "sudoer", flags.Migration.User)
		flagstests.AssertServerHelmFlags(t, &flags.Helm)
		flagstests.AssertVolumesFlags(t, &flags.Volumes)
		testutils.AssertEquals(t, "Error parsing --ssl-password", "sslsecret", flags.Installation.SSL.Password)
		testutils.AssertEquals(t, "Error parsing --ssh-key-public", "path/ssh.pub", flags.SSH.Key.Public)
		testutils.AssertEquals(t, "Error parsing --ssh-key-private", "path/ssh", flags.SSH.Key.Private)
		testutils.AssertEquals(t, "Error parsing --ssh-knownhosts", "path/known_hosts", flags.SSH.Knownhosts)
		testutils.AssertEquals(t, "Error parsing --ssh-config", "path/config", flags.SSH.Config)
		testutils.AssertEquals(t, "Wrong FQDN", "source.fq.dn", args[0])
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
