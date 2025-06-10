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
	args := []string{}

	args = append(args, "--ssl-password", "sslsecret")
	args = append(args, flagstests.ImageFlagsTestArgs...)
	args = append(args, flagstests.DBUpdateImageFlagTestArgs...)
	args = append(args, flagstests.CocoFlagsTestArgs...)
	args = append(args, flagstests.HubXmlrpcFlagsTestArgs...)
	args = append(args, flagstests.SalineFlagsTestArgs...)
	args = append(args, flagstests.PgsqlFlagsTestArgs...)
	args = append(args, flagstests.SCCFlagTestArgs...)
	args = append(args, flagstests.ServerKubernetesFlagsTestArgs...)
	args = append(args, flagstests.DBFlagsTestArgs...)
	args = append(args, flagstests.ReportDBFlagsTestArgs...)
	args = append(args, flagstests.InstallDBSSLFlagsTestArgs...)
	args = append(args, flagstests.SSLGenerationFlagsTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *kubernetes.KubernetesServerFlags,
		_ *cobra.Command, _ []string,
	) error {
		flagstests.AssertImageFlag(t, &flags.Image)
		flagstests.AssertDBUpgradeImageFlag(t, &flags.DBUpgradeImage)
		flagstests.AssertCocoFlag(t, &flags.Coco)
		flagstests.AssertHubXmlrpcFlag(t, &flags.HubXmlrpc)
		flagstests.AssertSalineFlag(t, &flags.Saline)
		flagstests.AssertPgsqlFlag(t, &flags.Pgsql)
		flagstests.AssertSCCFlag(t, &flags.Installation.SCC)
		flagstests.AssertServerKubernetesFlags(t, &flags.Kubernetes)
		flagstests.AssertDBFlag(t, &flags.Installation.DB)
		flagstests.AssertReportDBFlag(t, &flags.Installation.ReportDB)
		flagstests.AssertInstallDBSSLFlag(t, &flags.Installation.SSL.DB)
		flagstests.AssertSSLGenerationFlag(t, &flags.Installation.SSL.SSLCertGenerationFlags)
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
