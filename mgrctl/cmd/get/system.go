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

type System struct {
	ID                        int      `json:"id"`
	Name                      string   `json:"name"`
	LastCheckin               string   `json:"last_checkin"`
	Created                   string   `json:"created"`
	GroupCount                int      `json:"group_count"`
	SecurityErrata            int      `json:"security_errata"`
	BugErrata                 int      `json:"bug_errata"`
	EnhancementErrata         int      `json:"enhancement_errata"`
	OutdatedPkgCount          int      `json:"outdated_pkg_count"`
	ConfigFilesWithDifference int      `json:"config_files_with_difference"`
	ChannelLabels             string   `json:"channel_labels"`
	MgrServer                 bool     `json:"mgr_server"`
	Proxy                     bool     `json:"proxy"`
	Entitlement               []string `json:"entitlement"`
	VirtualHost               bool     `json:"virtual_host"`
	VirtualGuest              bool     `json:"virtual_guest"`
	ExtraPkgCount             int      `json:"extra_pkg_count"`
	RequiresReboot            bool     `json:"requires_reboot"`
	LastBoot                  string   `json:"last_boot"`
}

func init() {
	registerResource("system", &systemFetcher{}, []string{"sys"}, L("List systems"))
}

type systemFetcher struct{}

func (*systemFetcher) Columns() []utils.ColumnDef {
	return []utils.ColumnDef{
		{Header: "ID", Field: "ID"},
		{Header: "NAME", Field: "Name"},
		{Header: "LAST_CHECKIN", Field: "LastCheckin"},
		{Header: "CREATED", Field: "Created"},
	}
}

func (f *systemFetcher) List(client *api.APIClient, filter string, page, pageSize int) ([]System, int, error) {
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
	res, err := api.GetChecked[apitypes.FilteredResponse[System]](client, path, "api.system.list_systems_filtered")
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
