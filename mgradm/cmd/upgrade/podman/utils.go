// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func upgradePodman(globalFlags *types.GlobalFlags, flags *podmanUpgradeFlags, cmd *cobra.Command, args []string) error {
	return podman.Upgrade(flags.Image, flags.DbUpgradeImage, args)
}
