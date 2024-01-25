// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func InspectData(variable string, command string) types.Inspect {
	return types.Inspect{
		Variable: variable,
		Command:  command,
	}
}

var Values = []types.Inspect{
	InspectData("uyuni_release", "cat /etc/*release | grep 'Uyuni release' | cut -d ' ' -f3"),
	InspectData("suse_manager_release", "cat /etc/*release | grep 'SUSE Manager release' | cut -d ' ' -f4"),
	InspectData("image_pg_version", "rpm -qa --qf '%{VERSION}\n' 'name=postgresql[0-8][0-9]-server'  | cut -d. -f1 | sort -n | tail -1"),
	InspectData("current_pg_version", "cat /var/lib/pgsql/data/PG_VERSION"),
}

func AddInspectFlags(cmd *cobra.Command) {
	cmd_utils.AddImageFlag(cmd)
}

func GenerateInspectScript(scriptDir string) string {
	inspectScriptFilename := "inspect.sh"

	data := templates.InspectTemplateData{
		Param: values,
	}

	scriptPath := filepath.Join(scriptDir, inspectScriptFilename)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		log.Fatal().Err(err).Msgf("Failed to generate inspect script")
	}

	return inspectScriptFilename
}

func ReadInspectData(scriptDir string) map[string]string {
	data, err := os.ReadFile(filepath.Join(scriptDir, "data"))

	inspectResult := make(map[string]string)

	if err != nil {
		log.Fatal().Msgf("Failed to read data extracted from source host")
	}
	viper.SetConfigType("env")
	viper.ReadConfig(bytes.NewBuffer(data))

	for _, v := range values {
		if len(viper.GetString(v.Variable)) > 0 {
			inspectResult[v.Variable] = viper.GetString(v.Variable)
		}
	}
	return inspectResult
}
