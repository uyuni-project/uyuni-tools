// SPDX-FileCopyrightText: 2026 SUSE LLC
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
	hostData, err := podman.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman.PodmanLogin(hostData, flags.Image.Registry, flags.SCC)
	if err != nil {
		return err
	}

	defer cleaner()

	log.Debug().Msgf("Wanted database image %[1]s", flags.Pgsql.Image.Name)
	pgsqlImage, err := utils.ComputeImage("", utils.DefaultTag, flags.Pgsql.Image)
	if err != nil && len(pgsqlImage) > 0 {
		return utils.Errorf(err, L("failed to determine pgsql image"))
	}

	preparedServerImage, preparedDBImage, err := podman.PrepareImages(authFile, flags.Image, flags.Pgsql)
	if err != nil {
		return err
	}
	inspectResult, err := podman.Inspect(preparedServerImage, preparedDBImage)
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
