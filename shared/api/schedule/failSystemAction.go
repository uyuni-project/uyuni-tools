package schedule

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Fail specific event on specified system
func FailSystemAction(cnxDetails *api.ConnectionDetails, Sid int, ActionId int, Message string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"actionId":       ActionId,
		"message":       Message,
	}

	res, err := api.Post[types.#return_int_success()](client, "schedule", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute failSystemAction: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
