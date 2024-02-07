package schedule

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete all archived actions in the given list.
func DeleteActions(cnxDetails *api.ConnectionDetails, ActionIds []int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"actionIds":       ActionIds,
	}

	res, err := api.Post[types.#return_int_success()](client, "schedule", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteActions: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
