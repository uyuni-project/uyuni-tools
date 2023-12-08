// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/spf13/cobra"
)

var DefaultNamespace = "registry.opensuse.org/uyuni"
var DefaultTag = "latest"

// This variable needs to be set a build time using git tags
var Version = "0.0.0"

func AddBackendFlag(cmd *cobra.Command) {
	cmd.Flags().String("backend", "", "tool to use to reach the container. Possible values: 'podman', 'podman-remote', 'kubectl'. Default guesses which to use.")
}

// AddPullPolicyFlag adds the --pullPolicy flag to a command.
//
// Since podman doesn't have such a concept of pull policy like kubernetes,
// the values need some explanations for it:
//   - Never: just check and fail if needed
//   - IfNotPresent: check and pull
//   - Always: pull without checking
//
// For kubernetes the value is simply passed to the helm charts.
func AddPullPolicyFlag(cmd *cobra.Command) {
	cmd.Flags().String("pullPolicy", "IfNotPresent",
		"set whether to pull the images or not. The value can be one of 'Never', 'IfNotPresent' or 'Always'")
}
