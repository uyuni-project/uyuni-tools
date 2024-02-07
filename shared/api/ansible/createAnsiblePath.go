package ansible

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create ansible path
func CreateAnsiblePath(cnxDetails *api.ConnectionDetails, $param.getFlagName() $param.getType()) (*types.AnsiblePath, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.AnsiblePath](client, "ansible", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createAnsiblePath: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
