package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Transfer systems from one organization to another.  If executed by
 a #product() administrator, the systems will be transferred from their current
 organization to the organization specified by the toOrgId.  If executed by
 an organization administrator, the systems must exist in the same organization
 as that administrator and the systems will be transferred to the organization
 specified by the toOrgId. In any scenario, the origination and destination
 organizations must be defined in a trust.
func TransferSystems(cnxDetails *api.ConnectionDetails, ToOrgId int, Sids []int) (*types.#array_single("int", "serverIdTransferred"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"toOrgId":       ToOrgId,
		"sids":       Sids,
	}

	res, err := api.Post[types.#array_single("int", "serverIdTransferred")](client, "org", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute transferSystems: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
