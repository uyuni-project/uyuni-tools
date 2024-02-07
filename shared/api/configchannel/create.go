package configchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new global config channel. Caller must be at least a
 config admin or an organization admin.
func Create(cnxDetails *api.ConnectionDetails, Label string, Name string, Description string, Type string, $param.getFlagName() $param.getType()) (*types.ConfigChannel, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"name":       Name,
		"description":       Description,
		"type":       Type,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.ConfigChannel](client, "configchannel", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
