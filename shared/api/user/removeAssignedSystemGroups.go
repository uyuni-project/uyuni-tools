package user

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove system groups from a user's list of assigned system groups.
func RemoveAssignedSystemGroups(cnxDetails *api.ConnectionDetails, Login string, $param.getFlagName() $param.getType(), SetDefault bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"login":       Login,
		"$param.getName()":       $param.getFlagName(),
		"setDefault":       SetDefault,
	}

	res, err := api.Post[types.#return_int_success()](client, "user", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeAssignedSystemGroups: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
