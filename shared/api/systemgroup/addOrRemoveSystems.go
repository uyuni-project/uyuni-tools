package systemgroup

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Add/remove the given servers to a system group.
func AddOrRemoveSystems(cnxDetails *api.ConnectionDetails, SystemGroupName string, ServerIds []int, Add bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"systemGroupName":       SystemGroupName,
		"serverIds":       ServerIds,
		"add":       Add,
	}

	res, err := api.Post[types.#return_int_success()](client, "systemgroup", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addOrRemoveSystems: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
