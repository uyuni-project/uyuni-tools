// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// NewCommand for uninstall proxy.
func NewCommand(globalFlags *types.GlobalFlags) (*cobra.Command, error) {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: L("Uninstall a proxy"),
		Long: L(`Uninstall a proxy and optionally the corresponding volumes.
By default it will only print what would be done, use --force to actually remove.`) + kubernetes.UninstallHelp(),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			purge, _ := cmd.Flags().GetBool("purgeVolumes")

			backend, _ := cmd.Flags().GetString("backend")

			cnx := shared.NewConnection(backend, podman.ProxyContainerNames[0], kubernetes.ProxyFilter)
			command, err := cnx.GetCommand()
			if err != nil {
				return fmt.Errorf(L("failed to determine suitable backend: %s"), err)
			}
			switch command {
			case "podman":
				if err := uninstallForPodman(!force, purge); err != nil {
					return err
				}
			case "kubectl":
				if err := uninstallForKubernetes(!force); err != nil {
					return err
				}
			}
			return nil
		},
	}
	uninstallCmd.Flags().BoolP("force", "f", false, L("Actually remove the proxy"))
	uninstallCmd.Flags().Bool("purgeVolumes", false, L("Also remove the volumes"))

	utils.AddBackendFlag(uninstallCmd)

	return uninstallCmd, nil
}
