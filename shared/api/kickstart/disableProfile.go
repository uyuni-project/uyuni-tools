package kickstart

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Enable/Disable a Kickstart Profile
func DisableProfile(cnxDetails *api.ConnectionDetails, ProfileLabel string, Disabled string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"profileLabel":       ProfileLabel,
		"disabled":       Disabled,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute disableProfile: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
