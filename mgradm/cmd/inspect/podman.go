// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	inspect_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
)

const serverContainerName = "uyuni-server"

func InspectPodman(serverImage string, pullPolicy string) (map[string]string, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)

	if err != nil {
		return map[string]string{}, fmt.Errorf("Failed to create temporary directory %s", err)
	}

	extraArgs := []string{
		"-v", scriptDir + ":" + inspect_shared.InspectOutputFile.Directory,
	}

	shared_podman.PrepareImage(serverImage, pullPolicy)

	shared.GenerateInspectScript(scriptDir)

	podman.RunContainer("uyuni-inspect", serverImage, extraArgs,
		[]string{inspect_shared.InspectOutputFile.Directory + "/" + inspect_shared.InspectScriptFilename})

	inspectResult, err := shared.ReadInspectData(scriptDir)

	if err != nil {
		return map[string]string{}, fmt.Errorf("Cannot inspect data. %s", err)
	}

	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")

	if err != nil {
		return map[string]string{}, fmt.Errorf("Cannot print inspect result %s", err)
	}

	log.Info().Msgf("\n%s", string(prettyInspectOutput))
	return inspectResult, err

}
