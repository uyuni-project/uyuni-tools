package recurring

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete a recurring action with the given action ID.
func Delete(cnxDetails *api.ConnectionDetails, Id int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"id":       Id,
	}

	res, err := api.Post[types.#return_int_success()](client, "recurring", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute delete: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
