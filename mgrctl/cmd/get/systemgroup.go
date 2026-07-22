// SPDX-FileCopyrightText: 2026 Jayprakash
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type SystemGroup struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SystemCount string `json:"system_count"`
}

func init() {
	registerResource("systemgroup", systemGroupFetcher{}, []string{"grp"}, L("List system groups"))
}

type systemGroupFetcher struct{}

func (systemGroupFetcher) List(client *api.APIClient, _ string, _, _ int) ([]map[string]any, int, error) {
	res, err := api.GetChecked[[]map[string]any](client, "systemgroup/listAllGroups", "api.systemgroup.list_all_groups")
	if err != nil {
		return nil, 0, err
	}
	return res.Result, len(res.Result), nil
}

func (systemGroupFetcher) Columns() []utils.ColumnDef {
	return []utils.ColumnDef{
		{Header: "ID", Field: "ID"},
		{Header: "NAME", Field: "Name"},
		{Header: "DESCRIPTION", Field: "Description"},
		{Header: "SYSTEM_COUNT", Field: "SystemCount"},
	}
}
