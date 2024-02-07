package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set the manageable flag for a given channel and user.
 If value is set to 'true', this method will give the user
 manage permissions to the channel. Otherwise, that privilege is revoked.
func SetUserManageable(cnxDetails *api.ConnectionDetails, ChannelLabel string, Login string, Value bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
		"login":       Login,
		"value":       Value,
	}

	res, err := api.Post[types.#return_int_success()](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setUserManageable: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
