package ansible

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create ansible path
func RemoveAnsiblePath(cnxDetails *api.ConnectionDetails, PathId int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"pathId":       PathId,
	}

	res, err := api.Post[types.#return_int_success()](client, "ansible", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeAnsiblePath: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
