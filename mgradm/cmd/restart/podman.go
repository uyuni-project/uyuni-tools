// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package restart

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func podmanRestart(
	globalFlags *types.GlobalFlags,
	flags *restartFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return podman.RestartService(podman.ServerService)
}
