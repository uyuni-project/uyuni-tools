// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[kubernetes.KubernetesServerFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes [fqdn]",
		Short: L("Install a new server on a kubernetes cluster"),
		Long: L(`Install a new server on a kubernetes cluster

The install command assumes the following:
  * kubectl and helm are installed locally
  * a working kubectl configuration should be set to connect to the cluster to deploy to

The helm values file will be overridden with the values from the command parameters or configuration.

NOTE: installing on a remote cluster is not supported yet!
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetes.KubernetesServerFlags
			flags.ServerFlags.Coco.IsChanged = cmd.Flags().Changed("coco-replicas")
			flags.ServerFlags.HubXmlrpc.IsChanged = cmd.Flags().Changed("hubxmlrpc-replicas")
			return utils.CommandHelper(globalFlags, cmd, args, &flags, run)
		},
	}

	shared.AddInstallFlags(cmd)
	cmd_utils.AddHelmInstallFlag(cmd)
	cmd_utils.AddVolumesFlags(cmd)
	return cmd
}

// NewCommand for kubernetes installation.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, installForKubernetes)
}
