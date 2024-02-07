package ansible

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create ansible path
func UpdateAnsiblePath(cnxDetails *api.ConnectionDetails, PathId int, $param.getFlagName() $param.getType()) (*types.AnsiblePath, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"pathId":       PathId,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.AnsiblePath](client, "ansible", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updateAnsiblePath: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
