// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand for installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	installCmd := &cobra.Command{
		Use:     "install",
		GroupID: "deploy",
		Short:   L("Install a new server"),
		Long:    L("Install a new server"),
	}
	installCmd.PersistentFlags().StringVar(&globalFlags.Registry, "registry", "", L("specify a private registry"))

	installCmd.AddCommand(podman.NewCommand(globalFlags))

	if kubernetesCmd := kubernetes.NewCommand(globalFlags); kubernetesCmd != nil {
		installCmd.AddCommand(kubernetesCmd)
	}

	return installCmd
}
