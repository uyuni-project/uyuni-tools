package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update Content Project with given label
func UpdateProject(cnxDetails *api.ConnectionDetails, ProjectLabel string, $param.getFlagName() $param.getType()) (*types.ContentProject, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"projectLabel":       ProjectLabel,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.ContentProject](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updateProject: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
