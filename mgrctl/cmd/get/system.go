// SPDX-FileCopyrightText: 2026 Jayprakash
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	apitypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func init() {
	registerResource("system", systemFetcher{}, []string{"sys"}, L("List systems"))
}

type systemFetcher struct{}

func (systemFetcher) List(client *api.APIClient, filter string, page, pageSize int) ([]map[string]any, int, error) {
	filterKey, filterValue := "", ""
	if filter != "" {
		var err error
		filterKey, filterValue, err = parseFilter(filter)
		if err != nil {
			return nil, 0, err
		}
	}

	query := url.Values{}
	query.Set("filterKey", filterKey)
	query.Set("filterValue", filterValue)
	query.Set("page", fmt.Sprintf("%d", page))
	query.Set("pageSize", fmt.Sprintf("%d", pageSize))

	path := fmt.Sprintf("system/listSystemsFiltered?%s", query.Encode())
	res, err := api.GetChecked[apitypes.FilteredResponse[map[string]any]](client, path, "api.system.list_systems_filtered")
	if err != nil {
		return nil, 0, err
	}
	return res.Result.Data, res.Result.Total, nil
}

func parseFilter(expr string) (string, string, error) {
	var key, value string
	found := false
	for _, op := range []string{">=", "<=", "!=", "=", ">", "<"} {
		if i := strings.Index(expr, op); i >= 0 {
			key = strings.TrimSpace(expr[:i])
			value = strings.TrimSpace(expr[i:])
			found = true
			break
		}
	}

	if !found {
		key = strings.TrimSpace(expr)
	}

	if key == "" {
		return "", "", fmt.Errorf("filter key cannot be empty in %q", expr)
	}

	for _, c := range key {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' && c != '-' {
			return "", "", fmt.Errorf("invalid filter key %q: only letters, digits, underscores and dashes are allowed", key)
		}
	}

	return key, value, nil
}

func (systemFetcher) Columns() []utils.ColumnDef {
	return []utils.ColumnDef{
		{Header: "ID", Field: "id"},
		{Header: "NAME", Field: "name"},
		{Header: "LAST_CHECKIN", Field: "last_checkin"},
		{Header: "CREATED", Field: "created"},
	}
}
