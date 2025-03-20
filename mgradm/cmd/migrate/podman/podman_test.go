// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--prepare",
		"--user", "sudoer",
		"source.fq.dn",
	}

	args = append(args, flagstests.MirrorFlagTestArgs...)
	args = append(args, flagstests.SCCFlagTestArgs...)
	args = append(args, "--ssl-password", "sslsecret")
	args = append(args, flagstests.ImageFlagsTestArgs...)
	args = append(args, flagstests.DBUpdateImageFlagTestArgs...)
	args = append(args, flagstests.CocoFlagsTestArgs...)
	args = append(args, flagstests.HubXmlrpcFlagsTestArgs...)
	args = append(args, flagstests.SalineFlagsTestArgs...)
	args = append(args, flagstests.SCCFlagTestArgs...)
	args = append(args, flagstests.PodmanFlagsTestArgs...)
	args = append(args, flagstests.PgsqlFlagsTestArgs...)
	args = append(args, flagstests.DBFlagsTestArgs...)
	args = append(args, flagstests.ReportDBFlagsTestArgs...)
	args = append(args, flagstests.InstallDBSSLFlagsTestArgs...)
	args = append(args, flagstests.SSLGenerationFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *podmanMigrateFlags,
		_ *cobra.Command, args []string,
	) error {
		testutils.AssertTrue(t, "Prepare not set", flags.Migration.Prepare)
		flagstests.AssertMirrorFlag(t, flags.Mirror)
		testutils.AssertEquals(t, "Error parsing --user", "sudoer", flags.Migration.User)
		testutils.AssertEquals(t, "Wrong FQDN", "source.fq.dn", args[0])
		flagstests.AssertImageFlag(t, &flags.Image)
		flagstests.AssertDBUpgradeImageFlag(t, &flags.DBUpgradeImage)
		flagstests.AssertCocoFlag(t, &flags.Coco)
		flagstests.AssertHubXmlrpcFlag(t, &flags.HubXmlrpc)
		flagstests.AssertSalineFlag(t, &flags.Saline)
		flagstests.AssertSCCFlag(t, &flags.ServerFlags.Installation.SCC)
		flagstests.AssertPodmanInstallFlags(t, &flags.Podman)
		flagstests.AssertPgsqlFlag(t, &flags.Pgsql)
		flagstests.AssertDBFlag(t, &flags.Installation.DB)
		flagstests.AssertReportDBFlag(t, &flags.Installation.ReportDB)
		flagstests.AssertInstallDBSSLFlag(t, &flags.Installation.SSL.DB)
		flagstests.AssertSSLGenerationFlag(t, &flags.Installation.SSL.SSLCertGenerationFlags)
		testutils.AssertEquals(t, "Error parsing --ssl-password", "sslsecret", flags.ServerFlags.Installation.SSL.Password)
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
