package configchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove file paths from a global channel.
func DeleteFiles(cnxDetails *api.ConnectionDetails, Label string, Paths []string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"paths":       Paths,
	}

	res, err := api.Post[types.#return_int_success()](client, "configchannel", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteFiles: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
