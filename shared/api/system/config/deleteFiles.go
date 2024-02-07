package config

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Removes file paths from a local or sandbox channel of a server.
func DeleteFiles(cnxDetails *api.ConnectionDetails, Sid int, Paths []string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"paths":       Paths,
	}

	res, err := api.Post[types.#return_int_success()](client, "system/config", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteFiles: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
