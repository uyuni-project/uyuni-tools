package store

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List available image stores
func ListImageStores(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin() $ImageStoreSerializer #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "image/store"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin() $ImageStoreSerializer #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listImageStores: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
