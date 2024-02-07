package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new image profile
func Create(cnxDetails *api.ConnectionDetails, Label string, Type string, StoreLabel string, Path string, ActivationKey string, KiwiOptions string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"type":       Type,
		"storeLabel":       StoreLabel,
		"path":       Path,
		"activationKey":       ActivationKey,
		"kiwiOptions":       KiwiOptions,
	}

	res, err := api.Post[types.#return_int_success()](client, "image/profile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
