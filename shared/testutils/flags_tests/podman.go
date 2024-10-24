// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flags_tests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

// PodmanFlagsTestArgs is the values for PodmanFlagsTestArgs.
var PodmanFlagsTestArgs = []string{
	"--podman-arg", "arg1",
	"--podman-arg", "arg2",
}

// AssertPodmanInstallFlags checks that all podman flags are parsed correctly.
func AssertPodmanInstallFlags(t *testing.T, cmd *cobra.Command, flags *podman.PodmanFlags) {
	testutils.AssertEquals(t, "Error parsing --podman-arg", []string{"arg1", "arg2"}, flags.Args)
}
