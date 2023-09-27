// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
)

// This variable needs to be set a build time using git tags
var Version = "0.0.0"

func AddBackendFlag(cmd *cobra.Command) {
	cmd.Flags().String("backend", "", "tool to use to reach the container. Possible values: 'podman', 'podman-remote', 'kubectl'. Default guesses which to use.")
}

func AddAPIFlags(cmd *cobra.Command, apiFlags *api.ConnectionDetails) {
	cmd.PersistentFlags().String("host", "", "FQDN of the server to connect to")
	cmd.PersistentFlags().String("user", "", "API user username")
	cmd.PersistentFlags().String("password", "", "Password for the API user")

	// If host is not suplied, we try to take it from container using exec
	// The rest are mandatory
	cmd.MarkFlagRequired("user")
	cmd.MarkFlagRequired("password")
}
