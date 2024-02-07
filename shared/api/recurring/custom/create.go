package custom

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new recurring custom state action.
func Create(cnxDetails *api.ConnectionDetails, $param.getFlagName() $param.getType()) (*types.#param_desc("int", "id", "the ID of the newly created recurring action"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#param_desc("int", "id", "the ID of the newly created recurring action")](client, "recurring/custom", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
