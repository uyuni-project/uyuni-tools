package master

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Make this slave have no default Master for inter-server-sync
func UnsetDefaultMaster(cnxDetails *api.ConnectionDetails) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
	}

	res, err := api.Post[types.#return_int_success()](client, "sync/master", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute unsetDefaultMaster: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
