// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func upgradePodman(globalFlags *types.GlobalFlags, flags *podman.PodmanProxyFlags, cmd *cobra.Command, args []string) error {
	return podman.Upgrade(globalFlags, flags, cmd, args)
}
