package access

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set organization sharing access control.
func SetOrgSharing(cnxDetails *api.ConnectionDetails, ChannelLabel string, Access string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
		"access":       Access,
	}

	res, err := api.Post[types.#return_int_success()](client, "channel/access", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setOrgSharing: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
