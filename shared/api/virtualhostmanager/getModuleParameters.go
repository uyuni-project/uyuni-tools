package virtualhostmanager

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get a list of parameters for a virtual-host-gatherer module.
 It returns a map of parameters with their typical default values.
func GetModuleParameters(cnxDetails *api.ConnectionDetails, ModuleName string) (*types.#param_desc("map", "module_params", "module parameters"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "virtualhostmanager"
	params := ""
	if ModuleName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("map", "module_params", "module parameters")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getModuleParameters: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
