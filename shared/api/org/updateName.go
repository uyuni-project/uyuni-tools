package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Updates the name of an organization
func UpdateName(cnxDetails *api.ConnectionDetails, OrgId int, Name string) (*types.OrgDto, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"orgId": OrgId,
		"name":  Name,
	}

	res, err := api.Post[types.OrgDto](client, "org", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updateName: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
