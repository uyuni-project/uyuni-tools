// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"errors"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type SystemGroup struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OrgID       int    `json:"org_id"`
	SystemCount int    `json:"system_count"`
}

func init() {
	registerResource("systemgroup", &systemGroupFetcher{}, []string{"grp"}, false, false)
}

type systemGroupFetcher struct{}

func (*systemGroupFetcher) List(client *api.APIClient, _ ListOptions) ([]SystemGroup, int, error) {
	res, err := api.GetChecked[[]SystemGroup](client, "systemgroup/listAllGroups")
	if err != nil {
		return nil, 0, err
	}
	if !res.Success {
		return nil, 0, utils.Error(errors.New(res.Message), L("failed to get the system groups"))
	}

	return res.Result, len(res.Result), nil
}

func (*systemGroupFetcher) Columns() []utils.ColumnDef {
	return []utils.ColumnDef{
		{Header: "ID", Field: "ID"},
		{Header: "NAME", Field: "Name"},
		{Header: "DESCRIPTION", Field: "Description"},
		{Header: "ORG_ID", Field: "OrgID"},
		{Header: "SYSTEM_COUNT", Field: "SystemCount"},
	}
}

func (*systemGroupFetcher) Help() ResourceHelp {
	return ResourceHelp{
		Summary: L("List system groups"),
		Examples: L(`  # List all system groups
  mgrctl get systemgroup

  # Sort system groups by name
  mgrctl get systemgroup --sort-by=name`),
	}
}
