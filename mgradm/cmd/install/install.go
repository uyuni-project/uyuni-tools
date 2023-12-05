// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/podman"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	installCmd := &cobra.Command{
		Use:   "install [fqdn]",
		Short: "install a new server from scratch",
		Long:  "Install a new server from scratch",
	}

	installCmd.AddCommand(podman.NewCommand(globalFlags))

	if kubernetesCmd := kubernetes.NewCommand(globalFlags); kubernetesCmd != nil {
		installCmd.AddCommand(kubernetesCmd)
	}

	return installCmd
}
