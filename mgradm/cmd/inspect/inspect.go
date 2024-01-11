// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build !nok8s

package inspect

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"

	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type inspectFlags struct {
	Image      string
	Tag        string
	PullPolicy string
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	inspectCmd := &cobra.Command{
		Use:   "inspect",
		Short: "inspect",
		Long:  "Extract information from image and deployment",
		Args:  cobra.MaximumNArgs(0),

		RunE: func(cmd *cobra.Command, args []string) error {
			var flags inspectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, inspect)
		},
	}

	inspectCmd.Flags().String("image", "", "Image. Leave it empty to analyze the current deployment")
	inspectCmd.Flags().String("tag", "", "Tag Image. Leave it empty to analyze the current deployment")

	utils.AddPullPolicyFlag(inspectCmd)

	return inspectCmd
}

func inspect(globalFlags *types.GlobalFlags, flags *inspectFlags, cmd *cobra.Command, args []string) error {
	image, _ := cmd.Flags().GetString("image")
	tag, _ := cmd.Flags().GetString("tag")
	pullPolicy, _ := cmd.Flags().GetString("pullPolicy")

	backend := "podman"

	if kubernetesBuilt {
		backend, _ = cmd.Flags().GetString("backend")
	}

	cnx := shared.NewConnection(backend, serverContainerName, kubernetes.ServerFilter)
	command, err := cnx.GetCommand()

	if err != nil {
		return fmt.Errorf("Failed to determine suitable backend")
	}

	serverImage, err := utils.ComputeImage(image, tag)
	if err != nil && len(serverImage) > 0 {
		return fmt.Errorf("Failed to determine image. %s", err)
	}

	switch command {
	case "podman":
		if len(serverImage) <= 0 {
			log.Debug().Msg("Use deployed image")
			serverImage, err = adm_utils.RunningImage(cnx, serverContainerName)
			if err != nil {
				return fmt.Errorf("Failed to find current running image")
			}
		}
		_, err = InspectPodman(serverImage, pullPolicy)
	case "kubectl":
		if len(serverImage) <= 0 {
			log.Debug().Msg("Use deployed image")
			serverImage, err = adm_utils.RunningImage(cnx, "uyuni")
			if err != nil {
				return fmt.Errorf("Failed to find current running image. %s", err)
			}
		}
		_, err = InspectKubernetes(serverImage, pullPolicy)
	}
	return err
}
