package activationkey

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove server groups from an activation key.
func RemoveServerGroups(cnxDetails *api.ConnectionDetails, Key string, ServerGroupIds []int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"key":       Key,
		"serverGroupIds":       ServerGroupIds,
	}

	res, err := api.Post[types.#return_int_success()](client, "activationkey", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeServerGroups: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
