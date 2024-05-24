// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/install/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/install/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand install a new proxy from scratch.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install [fqdn]",
		Short: L("Install a new proxy from scratch"),
		Long:  L("Install a new proxy from scratch"),
	}
	installCmd.PersistentFlags().StringVar(&globalFlags.Registry, "registry", "", L("specify a private registry"))

	installCmd.AddCommand(podman.NewCommand(globalFlags))
	installCmd.AddCommand(kubernetes.NewCommand(globalFlags))

	return installCmd
}
