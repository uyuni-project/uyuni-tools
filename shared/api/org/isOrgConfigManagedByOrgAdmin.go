package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns whether Organization Administrator is able to manage his
 organization configuration. This may have a high impact on general #product() performance.
func IsOrgConfigManagedByOrgAdmin(cnxDetails *api.ConnectionDetails, OrgId int) (*types.#param_desc("boolean", "status", "Returns the status org admin management setting"), error) {
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

    res, err := api.Get[types.#param_desc("boolean", "status", "Returns the status org admin management setting")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute isOrgConfigManagedByOrgAdmin: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
