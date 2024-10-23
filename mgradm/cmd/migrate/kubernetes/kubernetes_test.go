// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils/flags_tests"
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

	args = append(args, flags_tests.MirrorFlagTestArgs...)
	args = append(args, flags_tests.SccFlagTestArgs...)
	args = append(args, flags_tests.ImageFlagsTestArgs...)
	args = append(args, flags_tests.DbUpdateImageFlagTestArgs...)
	args = append(args, flags_tests.CocoFlagsTestArgs...)
	args = append(args, flags_tests.HubXmlrpcFlagsTestArgs...)
	args = append(args, flags_tests.ServerHelmFlagsTestArgs...)
	args = append(args, flags_tests.VolumesFlagsTestExpected...)

	// Test function asserting that the args are properly parsed
	tester := func(globalFlags *types.GlobalFlags, flags *kubernetes.KubernetesServerFlags,
		cmd *cobra.Command, args []string,
	) error {
		test_utils.AssertTrue(t, "Prepare not set", flags.Migration.Prepare)
		flags_tests.AssertMirrorFlag(t, cmd, flags.Mirror)
		flags_tests.AssertSccFlag(t, cmd, &flags.Installation.Scc)
		flags_tests.AssertImageFlag(t, cmd, &flags.Image)
		flags_tests.AssertDbUpgradeImageFlag(t, cmd, &flags.DbUpgradeImage)
		flags_tests.AssertCocoFlag(t, cmd, &flags.Coco)
		flags_tests.AssertHubXmlrpcFlag(t, cmd, &flags.HubXmlrpc)
		test_utils.AssertEquals(t, "Error parsing --user", "sudoer", flags.Migration.User)
		flags_tests.AssertServerHelmFlags(t, cmd, &flags.Helm)
		flags_tests.AssertVolumesFlags(t, cmd, &flags.Volumes)
		test_utils.AssertEquals(t, "Error parsing --ssl-password", "sslsecret", flags.Installation.Ssl.Password)
		test_utils.AssertEquals(t, "Error parsing --ssh-key-public", "path/ssh.pub", flags.Ssh.Key.Public)
		test_utils.AssertEquals(t, "Error parsing --ssh-key-private", "path/ssh", flags.Ssh.Key.Private)
		test_utils.AssertEquals(t, "Error parsing --ssh-knownhosts", "path/known_hosts", flags.Ssh.Knownhosts)
		test_utils.AssertEquals(t, "Error parsing --ssh-config", "path/config", flags.Ssh.Config)
		test_utils.AssertEquals(t, "Wrong FQDN", "source.fq.dn", args[0])
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
