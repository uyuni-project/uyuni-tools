// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	apitypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
)

type systemResource struct{}

// List returns all systems. The API does not support filtering by name or ID, so the filtering is done client-side in the runGet function.
func (systemResource) List(client *api.APIClient) ([]apitypes.System, error) {
	res, err := api.Get[[]apitypes.System](client, "system/listSystems")
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

// Get by name or ID is not yet implemented because the API does not support it directly. We would need to fetch all systems and filter client-side, which is inefficient.
// This is a placeholder for future implementation when the API supports it.
func (systemResource) Get(client *api.APIClient, name string) (apitypes.System, error) {
	return apitypes.System{}, fmt.Errorf("get by name not yet implemented")
}

func (systemResource) Columns() []ColumnDef {
	return []ColumnDef{
		{Header: "ID", Field: "ID"},
		{Header: "NAME", Field: "Name"},
		{Header: "LAST_CHECKIN", Field: "LastCheckin"},
		{Header: "CREATED", Field: "Created"},
	}
}

func (systemResource) FilterFields() []string {
	return []string{"Name", "ID"}
}
