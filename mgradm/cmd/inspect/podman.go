// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	inspect_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
)

const serverContainerName = "uyuni-server"

func inspectPodman(serverImage string, pullPolicy string) (map[string]string, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)

	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create temporary directory")
	}

	extraArgs := []string{
		"-v", scriptDir + ":" + inspect_shared.InspectOutputFile.Directory,
	}

	shared_podman.PrepareImage(serverImage, pullPolicy)

	shared.GenerateInspectScript(scriptDir)

	podman.RunContainer("uyuni-inspect", serverImage, extraArgs,
		[]string{inspect_shared.InspectOutputFile.Directory + "/" + inspect_shared.InspectScriptFilename})

	inspectResult := shared.ReadInspectData(scriptDir)
	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")

	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot print inspect result")
	}

	log.Info().Msgf("\n%s", string(prettyInspectOutput))
	return inspectResult, err

}
