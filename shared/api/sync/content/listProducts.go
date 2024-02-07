package content

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all accessible products.
func ListProducts(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
                       $MgrSyncProductDtoSerializer
                    #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "sync/content"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
                       $MgrSyncProductDtoSerializer
                    #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listProducts: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
