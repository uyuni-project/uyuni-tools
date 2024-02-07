package systemgroup

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Add or remove administrators to/from the given group. #product() and
 Organization administrators are granted access to groups within their organization
 by default; therefore, users with those roles should not be included in the array
 provided. Caller must be an organization administrator.
func AddOrRemoveAdmins(cnxDetails *api.ConnectionDetails, SystemGroupName string, $param.getFlagName() $param.getType(), Add int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"systemGroupName":       SystemGroupName,
		"$param.getName()":       $param.getFlagName(),
		"add":       Add,
	}

	res, err := api.Post[types.#return_int_success()](client, "systemgroup", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addOrRemoveAdmins: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
