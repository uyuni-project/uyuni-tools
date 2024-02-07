package kickstart

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Import a kickstart profile.
func ImportFile(cnxDetails *api.ConnectionDetails, ProfileLabel string, VirtualizationType string, KickstartableTreeLabel string, KickstartFileContents string, KickstartHost string, UpdateType string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"profileLabel":       ProfileLabel,
		"virtualizationType":       VirtualizationType,
		"kickstartableTreeLabel":       KickstartableTreeLabel,
		"kickstartFileContents":       KickstartFileContents,
		"kickstartHost":       KickstartHost,
		"updateType":       UpdateType,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute importFile: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
