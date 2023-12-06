// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func GetConfigPath(args []string) string {
	configPath := args[0]
	if !utils.FileExists(configPath) {
		log.Fatal().Msgf("argument is not an existing file: %s", configPath)
	}
	return configPath
}
