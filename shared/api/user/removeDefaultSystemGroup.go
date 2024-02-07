package user

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove a system group from user's list of default system groups.
func RemoveDefaultSystemGroup(cnxDetails *api.ConnectionDetails, Login string, SgName string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"login":       Login,
		"sgName":       SgName,
	}

	res, err := api.Post[types.#return_int_success()](client, "user", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeDefaultSystemGroup: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
