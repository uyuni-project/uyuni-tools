package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove an action from an Action Chain.
func RemoveAction(cnxDetails *api.ConnectionDetails, ChainLabel string, ActionId int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"chainLabel":       ChainLabel,
		"actionId":       ActionId,
	}

	res, err := api.Post[types.#return_int_success()](client, "actionchain", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeAction: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
