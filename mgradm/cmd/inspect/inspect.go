// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

type inspectFlags struct {
	Image types.ImageFlags `mapstructure:",squash"`
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	inspectCmd := &cobra.Command{
		Use:   "inspect",
		Short: "inspect image",
		Long:  "Extract information from image",
	}

	inspectCmd.AddCommand(podman.NewCommand(globalFlags))

	//if kubernetesCmd := kubernetes.NewCommand(globalFlags); kubernetesCmd != nil {
	//	inspectCmd.AddCommand(kubernetesCmd)
	//}

	return inspectCmd
}
