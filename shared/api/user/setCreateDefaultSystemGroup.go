package user

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Sets the value of the createDefaultSystemGroup setting.
 If True this will cause there to be a system group created (with the same name
 as the user) every time a new user is created, with the user automatically given
 permission to that system group and the system group being set as the default
 group for the user (so every time the user registers a system it will be
 placed in that system group by default). This can be useful if different
 users will administer different groups of servers in the same organization.
 Can only be called by an org_admin.
func SetCreateDefaultSystemGroup(cnxDetails *api.ConnectionDetails, CreateDefaultSystemGroup bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"createDefaultSystemGroup":       CreateDefaultSystemGroup,
	}

	res, err := api.Post[types.#return_int_success()](client, "user", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setCreateDefaultSystemGroup: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
