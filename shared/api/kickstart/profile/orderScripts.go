package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Change the order that kickstart scripts will run for
 this kickstart profile. Scripts will run in the order they appear
 in the array. There are three arrays, one for all pre scripts, one
 for the post scripts that run before registration and server
 actions happen, and one for post scripts that run after registration
 and server actions. All scripts must be included in one of these
 lists, as appropriate.
func OrderScripts(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/profile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute orderScripts: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
