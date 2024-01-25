// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
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
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
	}

	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to compute image URL")
	}

	shared_podman.PrepareImage(serverImage, flags.Image.PullPolicy)

	inspectScriptBasename := GenerateInspectScript(scriptDir)

	podman.RunContainer("uyuni-inspect", serverImage, extraArgs,
		[]string{"/var/lib/uyuni-tools/" + inspectScriptBasename})

	inspectResult := readInspectData(scriptDir)
	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")

	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot print inspect result")
	}

	log.Info().Msgf("\n%s", string(prettyInspectOutput))
	return err

}

func GenerateInspectScript(scriptDir string) string {
	inspectScriptFilename := "inspect.sh"

	data := templates.InspectTemplateData{
		Param: shared.Values,
	}

	scriptPath := filepath.Join(scriptDir, inspectScriptFilename)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		log.Fatal().Err(err).Msgf("Failed to generate inspect script")
	}

	return inspectScriptFilename
}

func readInspectData(scriptDir string) map[string]string {
	data, err := os.ReadFile(filepath.Join(scriptDir, "data"))

	inspectResult := make(map[string]string)

	if err != nil {
		log.Fatal().Msgf("Failed to read data extracted from source host")
	}
	viper.SetConfigType("env")
	viper.ReadConfig(bytes.NewBuffer(data))

	for _, v := range shared.Values {
		inspectResult[v.Variable] = viper.GetString(v.Variable)
	}
	return inspectResult
}
