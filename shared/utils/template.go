// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"io"
	"os"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// Template is an interface for implementing Render function.
type Template interface {
	Render(wr io.Writer) error
}

// WriteTemplateToFile writes a template to a file.
func WriteTemplateToFile(template Template, path string, perm os.FileMode, overwrite bool) error {
	// Check if the file is existing
	if !overwrite {
		if FileExists(path) {
			return fmt.Errorf(L("%s file already present, not overwriting"), path)
		}
	}

	// Write the configuration
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return Errorf(err, L("failed to open %s for writing"), path)
	}
	defer file.Close()

	return template.Render(file)
}
