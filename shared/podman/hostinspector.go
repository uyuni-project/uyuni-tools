// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"

	"github.com/rs/zerolog"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// hostInspector inspects either the proxy or the server hosts.
type hostInspector struct {
	utils.BaseInspector
}

// newHostInspector creates a new HostInspector generating the inspection script and data in scriptDir.
func newHostInspector(scriptDir string) hostInspector {
	base := utils.BaseInspector{
		Values: []types.InspectData{
			types.NewInspectData(
				"scc_username",
				"cat /etc/zypp/credentials.d/SCCcredentials 2>&1 /dev/null | grep username | cut -d= -f2 || true"),
			types.NewInspectData(
				"scc_password",
				"cat /etc/zypp/credentials.d/SCCcredentials 2>&1 /dev/null | grep password | cut -d= -f2 || true"),
		},
		ScriptDir: scriptDir,
	}
	return hostInspector{
		BaseInspector: base,
	}
}

// hostInspectData are the data returned by the host inspector.
type hostInspectData struct {
	SccUsername string `mapstructure:"scc_username"`
	SccPassword string `mapstructure:"scc_password"`
}

// ReadInspectData parses the data generated by the host inspector.
func (i *hostInspector) ReadInspectData() (*hostInspectData, error) {
	return utils.ReadInspectData[hostInspectData](i.GetDataPath())
}

func inspectHost() (*hostInspectData, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return nil, utils.Errorf(err, L("failed to create temporary directory"))
	}

	inspector := newHostInspector(scriptDir)
	if err := inspector.GenerateScript(); err != nil {
		return nil, err
	}

	if err := utils.RunCmdStdMapping(zerolog.DebugLevel, inspector.GetScriptPath()); err != nil {
		return nil, utils.Errorf(err, L("failed to run inspect script in host system"))
	}

	inspectResult, err := inspector.ReadInspectData()
	if err != nil {
		return nil, utils.Errorf(err, L("cannot inspect host data"))
	}

	return inspectResult, err
}
