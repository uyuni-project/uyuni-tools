// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package podman

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func ptfForPodman(
	globalFlags *types.GlobalFlags,
	flags *podmanPTFFlags,
	cmd *cobra.Command,
	args []string,
) error {
	//we don't want to perform a postgres version upgrade when installing a PTF.
	//in that case, we can use the upgrade command.
	dummyMigration := types.ImageFlags{}
	return podman.Upgrade(flags.Image, dummyMigration, args)
}
