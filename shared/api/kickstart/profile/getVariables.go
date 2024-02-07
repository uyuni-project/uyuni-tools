package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of variables
                      associated with the specified kickstart profile
func GetVariables(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#struct_begin("kickstart variable")
         #prop("string", "key")
         #prop("string or int", "value")
     #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/profile"
	params := ""
	if KsLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("kickstart variable")
         #prop("string", "key")
         #prop("string or int", "value")
     #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getVariables: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
