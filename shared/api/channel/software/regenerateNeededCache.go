package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Completely clear and regenerate the needed Errata and Package
      cache for all systems subscribed to the specified channel.  This should
      be used only if you believe your cache is incorrect for all the systems
      in a given channel. This will schedule an asynchronous action to actually
      do the processing.
func RegenerateNeededCache(cnxDetails *api.ConnectionDetails, ChannelLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute regenerateNeededCache: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
