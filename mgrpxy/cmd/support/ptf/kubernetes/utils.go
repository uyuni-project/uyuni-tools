// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func ptfForKubernetes(globalFlags *types.GlobalFlags,
	flags *kubernetesPTFFlags,
	cmd *cobra.Command,
	args []string,
) error {

	return kubernetes.Upgrade(&flags.UpgradeFlags, cmd, args)
}
