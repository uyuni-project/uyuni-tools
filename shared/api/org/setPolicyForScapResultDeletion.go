package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set the status of SCAP result deletion settins for the given
 organization.
func SetPolicyForScapResultDeletion(cnxDetails *api.ConnectionDetails, OrgId int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"orgId":       OrgId,
	}

	res, err := api.Post[types.#return_int_success()](client, "org", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setPolicyForScapResultDeletion: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
