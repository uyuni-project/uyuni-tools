package config

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedules a deploy action for all the configuration files
 on the given list of systems.
func DeployAll(cnxDetails *api.ConnectionDetails, Date $type) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"date":       Date,
	}

	res, err := api.Post[types.#return_int_success()](client, "system/config", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deployAll: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
