// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package migrate

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	migrateCmd := &cobra.Command{
		Use:   "migrate [source server FQDN]",
		Short: "migrate a remote server to containers",
		Long:  "Migrate a remote server to containers",
	}

	migrateCmd.AddCommand(podman.NewCommand(globalFlags))

	if kubernetesCmd := kubernetes.NewCommand(globalFlags); kubernetesCmd != nil {
		migrateCmd.AddCommand(kubernetesCmd)
	}

	return migrateCmd
}
