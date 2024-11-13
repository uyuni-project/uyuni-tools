// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/shared"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type kubernetesMigrateFlags struct {
	shared.MigrateFlags `mapstructure:",squash"`
	Helm                cmd_utils.HelmFlags
	SCC                 types.SCCCredentials
	SSL                 cmd_utils.SSLCertFlags
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[kubernetesMigrateFlags]) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "kubernetes [source server FQDN]",
		Short: L("Migrate a remote server to containers running on a kubernetes cluster"),
		Long: L(`Migrate a remote server to containers running on a kubernetes cluster

This migration command assumes a few things:
  * the SSH configuration for the source server is complete, including user and
    all needed options to connect to the machine,
  * an SSH agent is started and the key to use to connect to the server is added to it,
  * kubectl and helm are installed locally,
  * a working kubectl configuration should be set to connect to the cluster to deploy to

When migrating a server with a automatically generated SSL Root CA certificate, the private key
password will be required to convert it to RSA in a kubernetes secret.
This is not needed if the source server does not have a generated SSL CA certificate.

NOTE: migrating to a remote cluster is not supported yet!
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetesMigrateFlags
			flagsUpdater := func(v *viper.Viper) {
				flags.MigrateFlags.Coco.IsChanged = v.IsSet("coco.replicas")
				flags.MigrateFlags.HubXmlrpc.IsChanged = v.IsSet("hubxmlrpc.replicas")
			}
			return utils.CommandHelper(globalFlags, cmd, args, &flags, flagsUpdater, run)
		},
	}

	shared.AddMigrateFlags(migrateCmd)
	cmd_utils.AddHelmInstallFlag(migrateCmd)
	migrateCmd.Flags().String("ssl-password", "", L("SSL CA generated private key password"))

	return migrateCmd
}

// NewCommand for kubernetes migration.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, migrateToKubernetes)
}
