package activationkey

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Given a list of activation keys and configuration channels,
 this method adds given configuration channels to either the top or
 the bottom (whichever you specify) of an activation key's
 configuration channels list. The ordering of the configuration channels
 provided in the add list is maintained while adding.
 If one of the configuration channels in the 'add' list
 already exists in an activation key, the
 configuration  channel will be re-ranked to the appropriate place.
func AddConfigChannels(cnxDetails *api.ConnectionDetails, Keys []string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"keys":       Keys,
	}

	res, err := api.Post[types.#return_int_success()](client, "activationkey", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addConfigChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
