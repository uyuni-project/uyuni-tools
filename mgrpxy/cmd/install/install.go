// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/install/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/install/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	installCmd := &cobra.Command{
		Use:   "install [fqdn]",
		Short: "install a new proxy from scratch",
		Long:  "Install a new proxy from scratch",
	}

	installCmd.AddCommand(podman.NewCommand(globalFlags))
	installCmd.AddCommand(kubernetes.NewCommand(globalFlags))

	return installCmd
}
