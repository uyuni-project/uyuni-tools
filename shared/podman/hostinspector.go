// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/rs/zerolog"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// NewHostInspector creates a new templates.InspectTemplateData for host data.
func NewHostInspector() templates.InspectTemplateData {
	return templates.InspectTemplateData{
		Values: []types.InspectData{
			types.NewInspectData(
				"scc_username",
				"cat /etc/zypp/credentials.d/SCCcredentials 2>/dev/null | grep username | cut -d= -f2 || true"),
			types.NewInspectData(
				"scc_password",
				"cat /etc/zypp/credentials.d/SCCcredentials 2>/dev/null | grep password | cut -d= -f2 || true"),

			types.NewInspectData(
				"has_uyuni_server",
				"systemctl list-unit-files uyuni-server.service >/dev/null && echo true || echo false"),
			types.NewInspectData(
				"has_salt_minion",
				"systemctl list-unit-files venv-salt-minion.service >/dev/null && echo true || echo false"),
		},
	}
}

// HostInspectData are the data returned by the host inspector.
type HostInspectData struct {
	SCCUsername    string `mapstructure:"scc_username"`
	SCCPassword    string `mapstructure:"scc_password"`
	HasUyuniServer bool   `mapstructure:"has_uyuni_server"`
	HasSaltMinion  bool   `mapstructure:"has_salt_minion"`
}

// InspectHost gathers data on the host where to install the server or proxy.
func InspectHost() (*HostInspectData, error) {
	inspector := NewHostInspector()
	script, err := inspector.GenerateScript()
	if err != nil {
		return nil, err
	}

	out, err := newRunner("bash", "-c", script).Log(zerolog.DebugLevel).Exec()
	if err != nil {
		return nil, utils.Errorf(err, L("failed to run inspect script in host system"))
	}

	inspectResult, err := utils.ReadInspectData[HostInspectData](out)
	if err != nil {
		return nil, utils.Errorf(err, L("cannot inspect host data"))
	}

	return inspectResult, err
}
