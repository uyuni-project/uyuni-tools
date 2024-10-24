// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func upgradeKubernetes(
	globalFlags *types.GlobalFlags,
	flags *kubernetesUpgradeFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return kubernetes.Upgrade(globalFlags, &flags.Image, &flags.DBUpgradeImage, flags.Helm, cmd, args)
}
