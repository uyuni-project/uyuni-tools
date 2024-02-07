package powermanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get current power management settings of the given system
func GetDetails(cnxDetails *api.ConnectionDetails, Sid int, Name string) (*types.#struct_begin("powerManagementParameters")
    #prop_desc("string", "powerType", "Power management type")
    #prop_desc("string", "powerAddress", "IP address for power management")
    #prop_desc("string", "powerUsername", "The Username")
    #prop_desc("string", "powerPassword", "The Password")
    #prop_desc("string", "powerId", "Identifier")
  #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system/provisioning/powermanagement"
	params := ""
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if Name {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("powerManagementParameters")
    #prop_desc("string", "powerType", "Power management type")
    #prop_desc("string", "powerAddress", "IP address for power management")
    #prop_desc("string", "powerUsername", "The Username")
    #prop_desc("string", "powerPassword", "The Password")
    #prop_desc("string", "powerId", "Identifier")
  #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
