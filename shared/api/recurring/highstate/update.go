package highstate

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update the properties of a recurring highstate action.
func Update(cnxDetails *api.ConnectionDetails, $param.getFlagName() $param.getType()) (*types.#param_desc("int", "id", "the ID of the updated recurring action"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#param_desc("int", "id", "the ID of the updated recurring action")](client, "recurring/highstate", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
