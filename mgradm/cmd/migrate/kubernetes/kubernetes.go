// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/shared"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type kubernetesMigrateFlags struct {
	shared.MigrateFlags `mapstructure:",squash"`
	Helm                cmd_utils.HelmFlags
	Ssl                 cmd_utils.SslCertFlags
}

// NewCommand for kubernetes migration.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "kubernetes [source server FQDN]",
		Short: "migrate a remote server to containers running on a kubernetes cluster",
		Long: `Migrate a remote server to containers running on a kubernetes cluster

This migration command assumes a few things:
  * the SSH configuration for the source server is complete, including user and
    all needed options to connect to the machine,
  * an SSH agent is started and the key to use to connect to the server is added to it,
  * kubectl is installed locally
  * A working kubeconfig should be set to connect to the cluster to deploy to

When migrating a server with a automatically generate SSL Root CA certificate, the private key
password will be required to convert it to RSA in order to be converted into a kubernetes secret.
This is not needed if the source server does not have a generated SSL CA certificate.

NOTE: for now installing on a remote cluster is not supported yet!
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetesMigrateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, migrateToKubernetes)
		},
	}

	shared.AddMigrateFlags(migrateCmd)
	cmd_utils.AddHelmInstallFlag(migrateCmd)
	migrateCmd.Flags().String("ssl-password", "", "SSL CA generated private key password")

	return migrateCmd
}
