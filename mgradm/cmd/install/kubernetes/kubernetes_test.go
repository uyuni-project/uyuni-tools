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
	args := flagstests.InstallFlagsTestArgs()
	args = append(args, flagstests.ServerKubernetesFlagsTestArgs...)
	args = append(args, flagstests.VolumesFlagsTestExpected...)
	args = append(args, flagstests.PgsqlFlagsTestArgs...)
	args = append(args, "srv.fq.dn")

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *kubernetes.KubernetesServerFlags,
		_ *cobra.Command, args []string,
	) error {
		flagstests.AssertInstallFlags(t, &flags.ServerFlags)
		flagstests.AssertServerKubernetesFlags(t, &flags.Kubernetes)
		flagstests.AssertVolumesFlags(t, &flags.Volumes)
		flagstests.AssertPgsqlFlag(t, &flags.Pgsql)
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
