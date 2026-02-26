// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/uyuni-project/uyuni-tools/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewServerInspector creates a new templates.InspectTemplateData for the big container inspection.
func NewServerInspector() templates.InspectTemplateData {
	return templates.InspectTemplateData{
		Values: []types.InspectData{
			types.NewInspectData(
				"uyuni_release",
				"cat /etc/*release | grep 'Uyuni release' | cut -d ' ' -f3 || true"),
			types.NewInspectData(
				"suse_manager_release",
				`[ -f /etc/susemanager-release ] && sed 's/.*(\([0-9.]\+\).*/\1/g' /etc/susemanager-release || true`),
			types.NewInspectData("libc_version", "ldd --version | head -n1 | sed 's/^ldd (GNU libc) //'"),
		},
	}
}

// ServerInspectData are the data extracted by a server inspector.
type ServerInspectData struct {
	UyuniRelease       string `mapstructure:"uyuni_release"`
	SuseManagerRelease string `mapstructure:"suse_manager_release"`
	LibcVersion        string `mapstructure:"libc_version"`
}

type InspectData struct {
	ServerInspectData    `mapstructure:",squash" json:"ServerImage"`
	DBInspectData        `mapstructure:",squash" json:"DBImage"`
	ContainerInspectData `mapstructure:",squash" json:"RunningContainer"`
}
