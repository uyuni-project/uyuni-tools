// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var inspectValues = []types.InspectData{
	types.InspectDataConstructor("uyuni_release", "cat /etc/*release | grep 'Uyuni release' | cut -d ' ' -f3"),
	types.InspectDataConstructor("suse_manager_release", "cat /etc/*release | grep 'SUSE Manager release' | cut -d ' ' -f4"),
	types.InspectDataConstructor("image_pg_version", "rpm -qa --qf '%{VERSION}\\n' 'name=postgresql[0-8][0-9]-server'  | cut -d. -f1 | sort -n | tail -1"),
	types.InspectDataConstructor("current_pg_version", "(test -e /var/lib/pgsql/data/PG_VERSION && cat /var/lib/pgsql/data/PG_VERSION) || true"),
}

var InspectOutputFile = types.InspectFile{
	Directory: "/var/lib/uyuni-tools",
	Basename:  "data",
}

var InspectScriptFilename = "inspect.sh"

func AddInspectFlags(cmd *cobra.Command) {
	cmd_utils.AddImageFlag(cmd)
}

func GenerateInspectScript(scriptDir string) {

	data := templates.InspectTemplateData{
		Param:      inspectValues,
		OutputFile: InspectOutputFile.Directory + "/" + InspectOutputFile.Basename,
	}

	scriptPath := filepath.Join(scriptDir, InspectScriptFilename)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		log.Fatal().Err(err).Msgf("Failed to generate inspect script")
	}
}

func ReadInspectData(scriptDir string) map[string]string {
	path := filepath.Join(scriptDir, "data")

	log.Debug().Msgf("Trying to read %s", path)

	data, err := os.ReadFile(path)

	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot parse file")
	}

	inspectResult := make(map[string]string)

	viper.SetConfigType("env")
	viper.ReadConfig(bytes.NewBuffer(data))

	for _, v := range inspectValues {
		if len(viper.GetString(v.Variable)) > 0 {
			inspectResult[v.Variable] = viper.GetString(v.Variable)
		}
	}
	return inspectResult
}
