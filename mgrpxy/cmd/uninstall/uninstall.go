// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall a proxy",
		Long:  "Uninstall a proxy and optionally the corresponding volumes." + kubernetes.UninstallHelp,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			purge, _ := cmd.Flags().GetBool("purge-volumes")

			backend := "podman"
			backend, _ = cmd.Flags().GetString("backend")

			cnx := shared.NewConnection(backend, podman.ProxyContainerNames[0], kubernetes.ProxyFilter)
			command, err := cnx.GetCommand()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to determine suitable backend")
			}
			switch command {
			case "podman":
				uninstallForPodman(dryRun, purge)
			case "kubectl":
				uninstallForKubernetes(dryRun)
			}
		},
	}
	uninstallCmd.Flags().BoolP("dry-run", "n", false, "Only show what would be done")
	uninstallCmd.Flags().Bool("purge-volumes", false, "Also remove the volume")

	utils.AddBackendFlag(uninstallCmd)

	return uninstallCmd
}
