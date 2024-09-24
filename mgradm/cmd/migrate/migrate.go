// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package migrate

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand for migration.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:     "migrate [source server FQDN]",
		GroupID: "deploy",
		Short:   L("Migrate a remote server to containers"),
		Long:    L("Migrate a remote server to containers"),
	}
	migrateCmd.AddCommand(podman.NewCommand(globalFlags))

	if kubernetesCmd := kubernetes.NewCommand(globalFlags); kubernetesCmd != nil {
		migrateCmd.AddCommand(kubernetesCmd)
	}

	return migrateCmd
}
