package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Adds an action to deploy a configuration file to an Action Chain.
func AddConfigurationDeployment(cnxDetails *api.ConnectionDetails, ChainLabel string, Sid int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"chainLabel":       ChainLabel,
		"sid":       Sid,
	}

	res, err := api.Post[types.#return_int_success()](client, "actionchain", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addConfigurationDeployment: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
