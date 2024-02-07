package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns whether errata e-mail notifications are enabled
 for the organization
func IsErrataEmailNotifsForOrg(cnxDetails *api.ConnectionDetails, OrgId int) (*types.#param_desc("boolean", "status", "Returns the status of the errata e-mail notification
 setting for the organization"), error) {
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

    res, err := api.Get[types.#param_desc("boolean", "status", "Returns the status of the errata e-mail notification
 setting for the organization")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute isErrataEmailNotifsForOrg: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
