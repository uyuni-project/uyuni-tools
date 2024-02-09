// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var inspectValues = []types.InspectData{
	types.NewInspectData("uyuni_release", "cat /etc/*release | grep 'Uyuni release' | cut -d ' ' -f3"),
	types.NewInspectData("suse_manager_release", "cat /etc/*release | grep 'SUSE Manager release' | cut -d ' ' -f4"),
	types.NewInspectData("fqdn", "cat /etc/rhn/rhn.conf | grep 'java.hostname' | cut -d' ' -f3"),
	types.NewInspectData("image_pg_version", "rpm -qa --qf '%{VERSION}\\n' 'name=postgresql[0-8][0-9]-server'  | cut -d. -f1 | sort -n | tail -1"),
	types.NewInspectData("current_pg_version", "(test -e /var/lib/pgsql/data/PG_VERSION && cat /var/lib/pgsql/data/PG_VERSION) || true"),
}

// InspectOutputFile represents the directory and the basename where the inspect values are stored.
var InspectOutputFile = types.InspectFile{
	Directory: "/var/lib/uyuni-tools",
	Basename:  "data",
}

// InspectScriptFilename is the inspect script basename.
var InspectScriptFilename = "inspect.sh"

// GenerateInspectScript create the inspect script.
func GenerateInspectScript(scriptDir string) error {
	data := templates.InspectTemplateData{
		Param:      inspectValues,
		OutputFile: InspectOutputFile.Directory + "/" + InspectOutputFile.Basename,
	}

	scriptPath := filepath.Join(scriptDir, InspectScriptFilename)
	if err := utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return fmt.Errorf("failed to generate inspect script: %s", err)
	}
	return nil
}

// ReadInspectData returns a map with the values inspected by an image and deploy.
func ReadInspectData(scriptDir string) (map[string]string, error) {
	path := filepath.Join(scriptDir, "data")
	log.Debug().Msgf("Trying to read %s", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]string{}, fmt.Errorf("cannot parse file %s: %s", path, err)
	}

	inspectResult := make(map[string]string)

	viper.SetConfigType("env")
	if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return map[string]string{}, fmt.Errorf("cannot read config: %s", err)
	}

	for _, v := range inspectValues {
		if len(viper.GetString(v.Variable)) > 0 {
			inspectResult[v.Variable] = viper.GetString(v.Variable)
		}
	}
	return inspectResult, nil
}
