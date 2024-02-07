package master

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Check if this host is reading configuration from an ISS master.
func HasMaster(cnxDetails *api.ConnectionDetails) (*types.#param_desc("boolean", "master", "True if has an ISS master, false otherwise"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
	}

	res, err := api.Post[types.#param_desc("boolean", "master", "True if has an ISS master, false otherwise")](client, "sync/master", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute hasMaster: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
