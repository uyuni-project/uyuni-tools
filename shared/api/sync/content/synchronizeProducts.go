package content

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Synchronize SUSE products between the Customer Center
             and the #product() database.
func SynchronizeProducts(cnxDetails *api.ConnectionDetails) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
	}

	res, err := api.Post[types.#return_int_success()](client, "sync/content", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute synchronizeProducts: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
