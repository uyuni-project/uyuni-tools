// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

type Template interface {
	Render(wr io.Writer) error
}

func WriteTemplateToFile(template Template, path string, perm os.FileMode, overwrite bool) error {
	// Check if the file is existing
	if !overwrite {
		if FileExists(path) {
			log.Fatal().Msgf("%s file already present, not overwriting", path)
		}
	}

	// Write the configuration
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to open %s for writing", path)
	}
	defer file.Close()

	return template.Render(file)
}
