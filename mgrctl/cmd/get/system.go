// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"fmt"
	"net/url"
	"strings"

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
	registerResource("system", &systemFetcher{}, []string{"sys"})
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

func (*systemFetcher) Help() ResourceHelp {
	return ResourceHelp{
		Summary: L("List systems"),
		Details: L("Systems support server-side filtering, sorting, and pagination."),
		Examples: L(`  # List systems that have extra packages installed
  mgrctl get system --filter "extra_pkg_count>0"

  # Show the first page of ten systems
  mgrctl get system --page 0 --page-size 10`),
	}
}

func (*systemFetcher) List(client *api.APIClient, opts ListOptions) ([]System, int, error) {
	filterValue := ""
	if opts.hasFilter() {
		// The listSystemsFiltered API takes the value without any operator (or simply !) for matches.
		if opts.FilterOperator == "~" || opts.FilterOperator == "!~" {
			filterValue = strings.Replace(opts.FilterOperator, "~", "", 1) + opts.FilterValue
		}
	}

	query := url.Values{}
	query.Set("filterKey", opts.FilterKey)
	query.Set("filterValue", filterValue)
	query.Set("page", fmt.Sprintf("%d", opts.PageNumber))
	query.Set("pageSize", fmt.Sprintf("%d", opts.PageSize))
	query.Set("sortKey", opts.SortBy)
	query.Set("sortDescending", fmt.Sprintf("%t", opts.SortDescending))

	path := fmt.Sprintf("system/listSystemsFiltered?%s", query.Encode())
	res, err := api.GetChecked[apitypes.FilteredResponse[System]](client, path)
	if err != nil {
		return nil, 0, err
	}
	return res.Result.Data, res.Result.Total, nil
}
