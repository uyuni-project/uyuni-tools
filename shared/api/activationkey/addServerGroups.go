package activationkey

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Add server groups to an activation key.
func AddServerGroups(cnxDetails *api.ConnectionDetails, Key string, ServerGroupIds []int) (*types.#return_int_success(), error) {
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
		return nil, fmt.Errorf("failed to execute addServerGroups: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
