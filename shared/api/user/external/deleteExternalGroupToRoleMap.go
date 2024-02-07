package external

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete the role map for an external group. Can only be called
 by a #product() Administrator.
func DeleteExternalGroupToRoleMap(cnxDetails *api.ConnectionDetails, Name string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
	}

	res, err := api.Post[types.#return_int_success()](client, "user/external", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteExternalGroupToRoleMap: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
