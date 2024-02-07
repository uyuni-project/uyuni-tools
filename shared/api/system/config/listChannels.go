package config

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all global('Normal', 'State') configuration channels associated to a
              system in the order of their ranking.
func ListChannels(cnxDetails *api.ConnectionDetails, Sid int) (*types.#return_array_begin()
  $ConfigChannelSerializer
  #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system/config"
	params := ""
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
  $ConfigChannelSerializer
  #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
