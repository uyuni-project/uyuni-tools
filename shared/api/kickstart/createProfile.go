package kickstart

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a kickstart profile.
func CreateProfile(cnxDetails *api.ConnectionDetails, ProfileLabel string, VirtualizationType string, KickstartableTreeLabel string, KickstartHost string, RootPassword string, UpdateType string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"profileLabel":       ProfileLabel,
		"virtualizationType":       VirtualizationType,
		"kickstartableTreeLabel":       KickstartableTreeLabel,
		"kickstartHost":       KickstartHost,
		"rootPassword":       RootPassword,
		"updateType":       UpdateType,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createProfile: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
