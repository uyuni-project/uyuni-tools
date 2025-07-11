// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func podmanInspect(
	_ *types.GlobalFlags,
	flags *inspectFlags,
	_ *cobra.Command,
	_ []string,
) error {
	serverImage, err := utils.ComputeImage(flags.Image)
	if err != nil && len(serverImage) > 0 {
		return utils.Errorf(err, L("failed to determine server image"))
	}

	if len(serverImage) <= 0 {
		log.Debug().Msg("Use already deployed server image")

		serverImage, err = podman.GetRunningImage(podman.ServerContainerName)
		flags.Image.Name = serverImage
		if err != nil {
			return utils.Errorf(err, L("failed to find the image of the currently running server container"))
		}
	}

	log.Debug().Msgf("Wanted database image %[1]s", flags.Pgsql.Image.Name)
	pgsqlImage, err := utils.ComputeImage(flags.Pgsql.Image)
	if err != nil && len(pgsqlImage) > 0 {
		return utils.Errorf(err, L("failed to determine pgsql image"))
	}

	if len(pgsqlImage) <= 0 {
		log.Debug().Msg("Use already deployed database image")

		pgsqlImage, err = podman.GetRunningImage(podman.DBContainerName)
		flags.Pgsql.Image.Name = pgsqlImage
		if err != nil {
			return utils.Errorf(err, L("failed to find the image of the currently running db container"))
		}
	}

	inspectResult, err := podman.Inspect(flags.Image, flags.Pgsql.Image, flags.SCC)
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
