// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flags_tests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
)

// Expected values for PodmanFlagsTestArgs.
var PodmanFlagsTestArgs = []string{
	"--podman-arg", "arg1",
	"--podman-arg", "arg2",
}

// Assert that all podman flags are parsed correctly.
func AssertPodmanInstallFlags(t *testing.T, cmd *cobra.Command, flags *podman.PodmanFlags) {
	test_utils.AssertEquals(t, "Error parsing --podman-arg", []string{"arg1", "arg2"}, flags.Args)
}
