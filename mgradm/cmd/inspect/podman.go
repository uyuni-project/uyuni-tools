// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	inspect_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
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
	serverImage, err := utils.ComputeImage(flags.Image, flags.Tag)
	if err != nil && len(serverImage) > 0 {
		return fmt.Errorf("failed to determine image. %s", err)
	}

	if len(serverImage) <= 0 {
		log.Debug().Msg("Use deployed image")

		cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")
		serverImage, err = adm_utils.RunningImage(cnx, shared_podman.ServerContainerName)
		if err != nil {
			return fmt.Errorf("failed to find current running image")
		}
	}
	inspectResult, err := InspectPodman(serverImage, flags.PullPolicy)
	if err != nil {
		return fmt.Errorf("inspect command failed %s", err)
	}
	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot print inspect result %s", err)
	}

	log.Info().Msgf("\n%s", string(prettyInspectOutput))

	return nil
}

// InspectPodman check values on a given image and deploy.
func InspectPodman(serverImage string, pullPolicy string) (map[string]string, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to create temporary directory %s", err)
	}

	extraArgs := []string{
		"-v", scriptDir + ":" + inspect_shared.InspectOutputFile.Directory,
	}

	err = shared_podman.PrepareImage(serverImage, pullPolicy)
	if err != nil {
		return map[string]string{}, err
	}

	if err := inspect_shared.GenerateInspectScript(scriptDir); err != nil {
		return map[string]string{}, err
	}

	err = podman.RunContainer("uyuni-inspect", serverImage, extraArgs,
		[]string{inspect_shared.InspectOutputFile.Directory + "/" + inspect_shared.InspectScriptFilename})
	if err != nil {
		return map[string]string{}, err
	}

	inspectResult, err := inspect_shared.ReadInspectData(scriptDir)
	if err != nil {
		return map[string]string{}, fmt.Errorf("cannot inspect data. %s", err)
	}

	return inspectResult, err
}
