package custominfo

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete an existing custom key and all systems' values for the key.
func DeleteKey(cnxDetails *api.ConnectionDetails, KeyLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"keyLabel":       KeyLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "system/custominfo", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteKey: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
