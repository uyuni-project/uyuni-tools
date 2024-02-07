package tree

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the available kickstartable trees for the given channel.
func List(cnxDetails *api.ConnectionDetails, ChannelLabel string) (*types.#return_array_begin() $KickstartTreeSerializer #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
	}

	res, err := api.Post[types.#return_array_begin() $KickstartTreeSerializer #array_end()](client, "kickstart/tree", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute list: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
