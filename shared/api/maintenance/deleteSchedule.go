package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove a maintenance schedule
func DeleteSchedule(cnxDetails *api.ConnectionDetails, Name string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
	}

	res, err := api.Post[types.#return_int_success()](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteSchedule: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
