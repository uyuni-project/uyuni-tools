package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Clone a list of errata into the specified channel.
func Clone(cnxDetails *api.ConnectionDetails, ChannelLabel string, $param.getFlagName() $param.getType()) (*types.#return_array_begin()
              $ErrataSerializer
          #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#return_array_begin()
              $ErrataSerializer
          #array_end()](client, "errata", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute clone: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
