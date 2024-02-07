package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update Content Environment with given label
func UpdateEnvironment(cnxDetails *api.ConnectionDetails, ProjectLabel string, EnvLabel string, $param.getFlagName() $param.getType()) (*types.ContentEnvironment, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"projectLabel":       ProjectLabel,
		"envLabel":       EnvLabel,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.ContentEnvironment](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updateEnvironment: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
