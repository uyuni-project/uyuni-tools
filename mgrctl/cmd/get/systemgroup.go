// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"errors"
	"fmt"
	"sort"
	"strings"

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

var systemGroupSortFields = []string{
	"id",
	"name",
	"description",
	"org_id",
	"system_count",
}

func systemGroupComparator(sortBy string) (func(a, b SystemGroup) bool, error) {
	switch sortBy {
	case "":
		return nil, nil
	case "id":
		return func(a, b SystemGroup) bool {
			return a.ID < b.ID
		}, nil
	case "name":
		return func(a, b SystemGroup) bool {
			return a.Name < b.Name
		}, nil
	case "description":
		return func(a, b SystemGroup) bool {
			return a.Description < b.Description
		}, nil
	case "org_id":
		return func(a, b SystemGroup) bool {
			return a.OrgID < b.OrgID
		}, nil
	case "system_count":
		return func(a, b SystemGroup) bool {
			return a.SystemCount < b.SystemCount
		}, nil
	default:
		return nil, fmt.Errorf(
			L("unsupported sort field %[1]q for system groups; supported fields: %[2]s"),
			sortBy,
			strings.Join(systemGroupSortFields, ", "),
		)
	}
}

func init() {
	registerResource[SystemGroup]("systemgroup", &systemGroupFetcher{}, []string{"grp"})
}

type systemGroupFetcher struct{}

func (*systemGroupFetcher) List(client *api.APIClient, opts ListOptions) ([]SystemGroup, int, error) {
	if opts.Filter != "" {
		return nil, 0, errors.New(L("filtering is not supported for system groups"))
	}

	compare, err := systemGroupComparator(opts.SortBy)
	if err != nil {
		return nil, 0, err
	}

	res, err := api.GetChecked[[]SystemGroup](client, "systemgroup/listAllGroups", "api.systemgroup.list_all_groups")
	if err != nil {
		return nil, 0, err
	}

	if compare != nil {
		sort.SliceStable(res.Result, func(i, j int) bool {
			return compare(res.Result[i], res.Result[j])
		})
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
		Details: fmt.Sprintf(
			L("System groups support client-side sorting by:\n  %s"),
			strings.Join(systemGroupSortFields, ", "),
		),
		Examples: L(`  # List all system groups
  mgrctl get systemgroup

  # Sort system groups by name
  mgrctl get systemgroup --sort-by=name`),
	}
}
