package admin

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of ssh connection data registered.
func List(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
     $PaygSshDataSerializer
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
	}

	res, err := api.Post[types.#return_array_begin()
     $PaygSshDataSerializer
 #array_end()](client, "admin", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute list: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
