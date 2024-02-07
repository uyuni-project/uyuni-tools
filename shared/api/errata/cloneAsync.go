package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Asynchronously clone a list of errata into the specified channel.
func CloneAsync(cnxDetails *api.ConnectionDetails, ChannelLabel string, $param.getFlagName() $param.getType()) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#return_int_success()](client, "errata", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute cloneAsync: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
