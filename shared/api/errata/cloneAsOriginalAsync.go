package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Asynchronously clones a list of errata into a specified cloned channel
 according the original erratas
func CloneAsOriginalAsync(cnxDetails *api.ConnectionDetails, ChannelLabel string, $param.getFlagName() $param.getType()) (*types.#return_int_success(), error) {
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
		return nil, fmt.Errorf("failed to execute cloneAsOriginalAsync: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
