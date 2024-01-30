// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type kubernetesInspectFlags struct {
	Image types.ImageFlags `mapstructure:",squash"`
	Helm  cmd_utils.HelmFlags
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	kubernetesCmd := &cobra.Command{
		Use:   "kubernetes",
		Short: "inspect kubernetes image",
		Long: `Extract information from kubernetes image and/or deployment
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetesInspectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, inspectKubernetes)
		},
	}

	shared.AddInspectFlags(kubernetesCmd)

	return kubernetesCmd
}
