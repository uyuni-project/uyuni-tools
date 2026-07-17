// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// variables for unit testing.
var inspectHost = podman.InspectHost
var podmanLogin = podman.PodmanLogin
var prepareImages = podman.PrepareImages
var inspectImages = podman.Inspect

func podmanInspect(
	_ *types.GlobalFlags,
	flags *inspectFlags,
	_ *cobra.Command,
	_ []string,
) error {
	hostData, err := inspectHost()
	if err != nil {
		return err
	}

	if !hostData.HasUyuniServer {
		return errors.New(L("server is not initialized."))
	}

	authFile, cleaner, err := podmanLogin(hostData, flags.Image.Registry, flags.SCC)
	if err != nil {
		return err
	}
	defer cleaner()

	preparedServerImage, preparedDBImage, err := prepareImages(authFile, flags.Image, flags.Pgsql)
	if err != nil {
		return err
	}
	inspectResult, err := inspectImages(preparedServerImage, preparedDBImage)
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
