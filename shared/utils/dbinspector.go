// SPDX-FileCopyrightText: 2026 SUSE LLC
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
			types.NewInspectData("db_image_pg_version",
				"echo $PG_MAJOR || true"),
			types.NewInspectData("db_image_libc_version", "ldd --version | head -n1 | sed 's/^ldd (GNU libc) //'"),
		},
	}
}

// DBInspectData are data of the DB data.
type DBInspectData struct {
	PgVersion   string `mapstructure:"db_image_pg_version"`
	LibcVersion string `mapstructure:"db_image_libc_version"`
}
