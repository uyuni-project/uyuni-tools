package activationkey

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List configuration channels
 associated to an activation key.
func ListConfigChannels(cnxDetails *api.ConnectionDetails, Key string) (*types.#return_array_begin()
     $ConfigChannelSerializer
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "activationkey"
	params := ""
	if Key {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
     $ConfigChannelSerializer
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listConfigChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
