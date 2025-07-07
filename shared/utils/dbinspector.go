// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/uyuni-project/uyuni-tools/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewDBInspector creates a new templates.InspectTemplateData for the database info.
func NewDBInspector() templates.InspectTemplateData {
	return templates.InspectTemplateData{
		Values: []types.InspectData{
			types.NewInspectData("image_pg_version",
				"echo $PG_MAJOR || true"),
			types.NewInspectData("image_libc_version", "ldd --version | head -n1 | sed 's/^ldd (GNU libc) //'"),
		},
	}
}

// DBInspectData are data of the DB data.
type DBInspectData struct {
	ImagePgVersion   string `mapstructure:"image_pg_version"`
	ImageLibcVersion string `mapstructure:"image_libc_version"`
}
