// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func InspectData(variable string, command string) types.Inspect {
	return types.Inspect{
		Variable: variable,
		Command:  command,
	}
}

var Values = []types.Inspect{
	InspectData("uyuni_release", "cat /etc/*release | grep 'Uyuni release' | cut -d ' ' -f3"),
	InspectData("suma_release", "cat /etc/*release | grep 'SUSE Manager release' | cut -d ' ' -f4"),
	InspectData("new_pg_version", "rpm -qa --qf '%{VERSION}\n' 'name=postgresql[0-8][0-9]-server'  | cut -d. -f1 | sort -n | tail -1"),
	//InspectData("old_pg_version", "cat /var/lib/pgsql/data/PG_VERSION"),
	//InspectData("Timezone", "timedatectl show -p Timezone"),
}

// The port names should be less than 15 characters long and lowercased for traefik to eat them
func inspectImage(
	globalFlags *types.GlobalFlags,
	flags *inspectFlags,
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
	log.Info().Msgf("%s", inspectResult)
	return err

}

func GenerateInspectScript(scriptDir string) string {
	inspectScriptFilename := "inspect.sh"

	data := templates.InspectTemplateData{
		Param: Values,
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

	for _, v := range Values {
		inspectResult[v.Variable] = viper.GetString(v.Variable)
	}
	return inspectResult
}
