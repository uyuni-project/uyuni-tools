// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package scripts

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const extractInfoScript = `#!/bin/bash
echo "Extracting Uyuni version..."
echo "uyuni_release=$(cat /etc/*release | grep 'Uyuni release' | cut -d ' ' -f3)" >> /var/lib/uyuni-tools/data
echo "Extracting SUSE Manager version..."
echo "suma_release=$(cat /etc/*release | grep 'SUSE Manager release' | cut -d ' ' -f4)" >> /var/lib/uyuni-tools/data
echo "Extracting postgresql versions..."
echo "new_pg_version=$(rpm -qa --qf '%{VERSION}\n' 'name=postgresql[0-8][0-9]-server'  | cut -d. -f1 | sort -n | tail -1)" >> /var/lib/uyuni-tools/data
echo "old_pg_version=$(cat /var/lib/pgsql/data/PG_VERSION)" >> /var/lib/uyuni-tools/data
echo "DONE"`

func GenerateExtractInfoScript(dir string, filename string) string {
	path := filepath.Join(dir, filename)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot create %s", path)
		return ""
	}
	_, err = f.WriteString(extractInfoScript)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot write in  %s", path)
		return ""
	}
	err = os.Chmod(path, 0755)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot set permission to %s", path)
		return ""
	}
	log.Debug().Msgf("New script generated at %s", path)
	return path
}
