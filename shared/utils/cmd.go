// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/spf13/cobra"
)

// This variable needs to be set a build time using git tags
var Version = "0.0.0"

func AddBackendFlag(cmd *cobra.Command) {
	cmd.Flags().String("backend", "", "tool to use to reach the container. Possible values: 'podman', 'podman-remote', 'kubectl'. Default guesses which to use.")
}
