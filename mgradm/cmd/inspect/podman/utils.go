// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	inspect_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// The port names should be less than 15 characters long and lowercased for traefik to eat them
func inspectImagePodman(
	globalFlags *types.GlobalFlags,
	flags *podmanInspectFlags,
	cmd *cobra.Command,
	args []string,
) error {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)

	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create temporary directory")
	}

	extraArgs := []string{
		"-v", scriptDir + ":" + inspect_shared.InspectOutputFile.Directory,
	}

	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to compute image URL")
	}

	shared_podman.PrepareImage(serverImage, flags.Image.PullPolicy)

	shared.GenerateInspectScript(scriptDir)

	podman.RunContainer("uyuni-inspect", serverImage, extraArgs,
		[]string{inspect_shared.InspectOutputFile.Directory + "/" + inspect_shared.InspectScriptFilename})

	inspectResult := shared.ReadInspectData(scriptDir)
	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")

	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot print inspect result")
	}

	log.Info().Msgf("\n%s", string(prettyInspectOutput))
	return err

}
