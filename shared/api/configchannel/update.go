package configchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update a global config channel. Caller must be at least a
 config admin or an organization admin, or have access to a system containing this
 config channel.
func Update(cnxDetails *api.ConnectionDetails, Label string, Name string, Description string) (*types.ConfigChannel, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"name":       Name,
		"description":       Description,
	}

	res, err := api.Post[types.ConfigChannel](client, "configchannel", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
