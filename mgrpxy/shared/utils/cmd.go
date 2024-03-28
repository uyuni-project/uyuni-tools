// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// GetConfigPath returns the configuration path if exists.
func GetConfigPath(args []string) string {
	configPath := args[0]
	if !utils.FileExists(configPath) {
		log.Fatal().Msgf(L("argument is not an existing file: %s"), configPath)
	}
	return configPath
}
