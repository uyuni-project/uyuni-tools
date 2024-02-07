package access

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Enable user restrictions for the given channel. If enabled, only
 selected users within the organization may subscribe to the channel.
func EnableUserRestrictions(cnxDetails *api.ConnectionDetails, ChannelLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "channel/access", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute enableUserRestrictions: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
