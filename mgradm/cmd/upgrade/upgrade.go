// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
package upgrade

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand for upgrading a local server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade server",
		Short: L("Upgrade local server"),
		Long:  L("Upgrade local server"),
	}

	upgradeCmd.AddCommand(podman.NewCommand(globalFlags))

	if kubernetesCmd := kubernetes.NewCommand(globalFlags); kubernetesCmd != nil {
		upgradeCmd.AddCommand(kubernetesCmd)
	}

	return upgradeCmd
}
