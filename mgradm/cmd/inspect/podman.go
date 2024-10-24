// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func podmanInspect(
	globalFlags *types.GlobalFlags,
	flags *inspectFlags,
	cmd *cobra.Command,
	args []string,
) error {
	serverImage, err := utils.ComputeImage("", utils.DefaultTag, flags.Image)
	if err != nil && len(serverImage) > 0 {
		return utils.Errorf(err, L("failed to determine image"))
	}

	if len(serverImage) <= 0 {
		log.Debug().Msg("Use deployed image")

		cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")
		serverImage, err = adm_utils.RunningImage(cnx, shared_podman.ServerContainerName)
		if err != nil {
			return utils.Errorf(err, L("failed to find the image of the currently running server container"))
		}
	}
	inspectResult, err := shared_podman.Inspect(serverImage, flags.Image.PullPolicy, flags.SCC, false)
	if err != nil {
		return utils.Errorf(err, L("inspect command failed"))
	}
	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")
	if err != nil {
		return utils.Errorf(err, L("cannot print inspect result"))
	}

	outputString := "\n" + string(prettyInspectOutput)
	log.Info().Msgf(outputString)

	return nil
}
