package config

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Given a list of servers and configuration channels,
 this method appends the configuration channels to either the top or
 the bottom (whichever you specify) of a system's subscribed
 configuration channels list. The ordering of the configuration channels
 provided in the add list is maintained while adding.
 If one of the configuration channels in the 'add' list
 has been previously subscribed by a server, the
 subscribed channel will be re-ranked to the appropriate place.
func AddChannels(cnxDetails *api.ConnectionDetails) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
	}

	res, err := api.Post[types.#return_int_success()](client, "system/config", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
