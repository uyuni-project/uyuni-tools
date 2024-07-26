// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package scale

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type scaleFlags struct {
	Backend  string
	Replicas int
}

func addScaleFlags(cmd *cobra.Command) {
	cmd.Flags().Int("replicas", 0, L("How many replicas of a service should be started."))
}

// NewCommand adjusts a containers replicas.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	scaleCmd := &cobra.Command{
		Use:     "scale",
		GroupID: "management",
		Short:   L("Adjust the replicas for services supporting it."),
		Long: L(`Adjust the replicas for services supporting it.
Supported services:
  - uyuni-server-attestation
  - uyuni-hub-xmlrpc
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scaleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scale)
		},
	}
	scaleCmd.SetUsageTemplate(scaleCmd.UsageTemplate())
	addScaleFlags(scaleCmd)

	if utils.KubernetesBuilt {
		utils.AddBackendFlag(scaleCmd)
	}

	return scaleCmd
}

func scale(globalFlags *types.GlobalFlags, flags *scaleFlags, cmd *cobra.Command, args []string) error {
	fn, err := shared.ChoosePodmanOrKubernetes(cmd.Flags(), podmanScale, kubernetesScale)
	if err != nil {
		return err
	}

	return fn(globalFlags, flags, cmd, args)
}
