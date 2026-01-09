// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[kubernetes.KubernetesServerFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes [source server FQDN]",
		Short: L("Migrate a remote server to containers running on a kubernetes cluster"),
		Long: L(`Migrate a remote server to containers running on a kubernetes cluster

This migration command assumes a few things:
  * the SSH configuration for the source server is complete, including user and
    all needed options to connect to the machine,
  * kubectl and helm are installed locally,
  * a working kubectl configuration should be set to connect to the cluster to deploy to

The SSH parameters may be left empty if the target Kubernetes namespace contains:
  * an uyuni-migration-config ConfigMap with config and known_hosts items,
  * an uyuni-migration-key secret with key and key.pub items with a passwordless key.

When migrating a server with a automatically generated SSL Root CA certificate, the private key
password will be required to convert it to RSA in a kubernetes secret.
This is not needed if the source server does not have a generated SSL CA certificate.
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetes.KubernetesServerFlags
			flagsUpdater := func(v *viper.Viper) {
				flags.Coco.IsChanged = v.IsSet("coco.replicas")
				flags.HubXmlrpc.IsChanged = v.IsSet("hubxmlrpc.replicas")
				flags.Saline.IsChanged = v.IsSet("saline.replicas") || v.IsSet("saline.port")
				flags.Pgsql.IsChanged = v.IsSet("pgsql.replicas")
			}
			return utils.CommandHelper(globalFlags, cmd, args, &flags, flagsUpdater, run)
		},
	}
	cmd_utils.AddMirrorFlag(cmd)
	shared.AddMigrateFlags(cmd)
	cmd_utils.AddHelmInstallFlag(cmd)
	cmd_utils.AddVolumesFlags(cmd)

	cmd.Flags().String("ssh-key-public", "", L("Path to the SSH public key to use to connect to the source server"))
	cmd.Flags().String("ssh-key-private", "",
		L("Path to the passwordless SSH private key to use to connect to the source server"),
	)
	cmd.Flags().String("ssh-knownhosts", "", L("Path to the SSH known_hosts file to use to connect to the source server"))
	cmd.Flags().String("ssh-config", "", L("Path to the SSH configuration file to use to connect to the source server"))

	const sshGroupID = "ssh"
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: sshGroupID, Title: L("SSH Configuration Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "ssh-key-public", sshGroupID)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssh-key-private", sshGroupID)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssh-knownhosts", sshGroupID)
	_ = utils.AddFlagToHelpGroupID(cmd, "ssh-config", sshGroupID)

	return cmd
}

// NewCommand for kubernetes migration.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, migrateToKubernetes)
}
