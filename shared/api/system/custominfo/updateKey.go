package custominfo

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update description of a custom key
func UpdateKey(cnxDetails *api.ConnectionDetails, KeyLabel string, KeyDescription string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"keyLabel":       KeyLabel,
		"keyDescription":       KeyDescription,
	}

	res, err := api.Post[types.#return_int_success()](client, "system/custominfo", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updateKey: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
