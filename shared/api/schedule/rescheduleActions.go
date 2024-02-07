package schedule

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Reschedule all actions in the given list.
func RescheduleActions(cnxDetails *api.ConnectionDetails, ActionIds []int, OnlyFailed bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"actionIds":       ActionIds,
		"onlyFailed":       OnlyFailed,
	}

	res, err := api.Post[types.#return_int_success()](client, "schedule", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rescheduleActions: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
