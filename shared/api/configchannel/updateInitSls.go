package configchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update the init.sls file for the given state channel. User can only update contents, nothing else.
func UpdateInitSls(cnxDetails *api.ConnectionDetails, Label string, $param.getFlagName() $param.getType()) (*types.ConfigRevision, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.ConfigRevision](client, "configchannel", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updateInitSls: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
