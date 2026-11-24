// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package scale

import (
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type scaleFlags struct {
	Replicas  int
	Image     types.ImageFlags     `mapstructure:",squash"`
	SCC       types.SCCCredentials `mapstructure:"scc"`
	Coco      adm_utils.CocoFlags
	HubXmlrpc adm_utils.HubXmlrpcFlags
	Saline    adm_utils.SalineFlags
}

func addScaleFlags(cmd *cobra.Command) {
	cmd.Flags().Int("replicas", 0, L("How many replicas of a service should be started."))
	utils.AddPullPolicyFlag(cmd)
	utils.AddRegistryFlag(cmd)
	adm_utils.AddSCCFlag(cmd)

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "coco-container", Title: L("Confidential Computing Flags")})
	utils.AddContainerImageFlags(cmd, "coco", L("Confidential computing attestation"),
		"coco-container", "server-attestation")
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "hubxmlrpc-container", Title: L("Hub XML-RPC API")})
	utils.AddContainerImageFlags(cmd, "hubxmlrpc", L("Hub XML-RPC API"), "hubxmlrpc-container", "server-hub-xmlrpc-api")
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "saline-container", Title: L("Saline Flags")})
	utils.AddContainerImageFlags(cmd, "saline", L("Saline"), "saline-container", "server-saline")
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[scaleFlags]) *cobra.Command {
	scaleCmd := &cobra.Command{
		Use:     "scale",
		GroupID: "management",
		Short:   L("Adjust the replicas for services supporting it."),
		Long: L(`Adjust the replicas for services supporting it.
Supported services:
  - uyuni-hub-xmlrpc
  - uyuni-saline
  - uyuni-server-attestation

The current image of the service is kept, unless its image or tag flag is passed.
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scaleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	scaleCmd.SetUsageTemplate(scaleCmd.UsageTemplate())
	addScaleFlags(scaleCmd)

	return scaleCmd
}

// NewCommand adjusts a containers replicas.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, podmanScale)
}
