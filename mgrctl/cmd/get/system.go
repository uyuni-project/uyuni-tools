// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"errors"
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
	GroupCount                int      `json:"group_count"`
	SecurityErrata            int      `json:"security_errata"`
	BugErrata                 int      `json:"bug_errata"`
	EnhancementErrata         int      `json:"enhancement_errata"`
	OutdatedPackages          int      `json:"outdated_packages"`
	TotalErrataCount          int      `json:"total_errata_count"`
	ConfigFilesWithDifference int      `json:"config_files_with_difference"`
	ChannelLabels             string   `json:"channel_labels"`
	LastCheckin               string   `json:"last_checkin"`
	LastCheckinDays           int      `json:"last_checkin_days"`
	MgrServer                 bool     `json:"mgr_server"`
	Proxy                     bool     `json:"proxy"`
	Entitlement               []string `json:"entitlement"`
	VirtualHost               bool     `json:"virtual_host"`
	VirtualGuest              bool     `json:"virtual_guest"`
	ExtraPkgCount             int      `json:"extra_pkg_count"`
	RequiresReboot            bool     `json:"requires_reboot"`
	ActionsCount              int      `json:"actions_count"`
	Kickstarting              bool     `json:"kickstarting"`
	UnscheduledErrataCount    int      `json:"unscheduled_errata_count"`
	CreatorName               string   `json:"creator_name"`
	Created                   string   `json:"created"`
	CreatedDays               int      `json:"created_days"`
	StatusType                string   `json:"status_type"`
}

func init() {
	registerResource("system", &systemFetcher{}, []string{"sys"}, true, true)
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
		Examples: L(`  # List systems that have extra packages installed
  mgrctl get system --filter "extra_pkg_count>0"

  # Show the first page of ten systems
  mgrctl get system --page 1 --page-size 10`),
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
	if !res.Success {
		return nil, 0, utils.Error(errors.New(res.Message), L("failed to get the systems"))
	}
	return res.Result.Data, res.Result.Total, nil
}
