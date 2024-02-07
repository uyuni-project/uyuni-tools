package external

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update the roles for an external group. Replace previously set roles
 with the ones passed in here. Can only be called by a #product() Administrator.
func SetExternalGroupRoles(cnxDetails *api.ConnectionDetails, Name string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
	}

	res, err := api.Post[types.#return_int_success()](client, "user/external", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setExternalGroupRoles: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
