package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Sets whether Organization Administrator can manage his organization
 configuration. This may have a high impact on general #product() performance.
func SetOrgConfigManagedByOrgAdmin(cnxDetails *api.ConnectionDetails, OrgId int, Enable bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"orgId":       OrgId,
		"enable":       Enable,
	}

	res, err := api.Post[types.#return_int_success()](client, "org", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setOrgConfigManagedByOrgAdmin: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
