package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Clones a list of errata into a specified cloned channel according the original erratas.
func CloneAsOriginal(cnxDetails *api.ConnectionDetails, ChannelLabel string, $param.getFlagName() $param.getType()) (*types.#return_array_begin()
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
		return nil, fmt.Errorf("failed to execute cloneAsOriginal: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
