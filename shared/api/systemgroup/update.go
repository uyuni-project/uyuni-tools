package systemgroup

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update an existing system group.
func Update(cnxDetails *api.ConnectionDetails, SystemGroupName string, Description string) (*types.ManagedServerGroup, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"systemGroupName":       SystemGroupName,
		"description":       Description,
	}

	res, err := api.Post[types.ManagedServerGroup](client, "systemgroup", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
