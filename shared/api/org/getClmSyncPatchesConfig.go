package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Reads the content lifecycle management patch synchronization config option.
func GetClmSyncPatchesConfig(cnxDetails *api.ConnectionDetails, OrgId int) (*types.#param_desc("boolean", "status", "Get the config option value"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "org"
	params := ""
	if OrgId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("boolean", "status", "Get the config option value")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getClmSyncPatchesConfig: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
