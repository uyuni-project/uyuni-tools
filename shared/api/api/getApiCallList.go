package api

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lists all available api calls grouped by namespace
func GetApiCallList(cnxDetails *api.ConnectionDetails) (*types.#struct_begin("method_info")
       #prop_desc("string", "name", "method name")
       #prop_desc("string", "parameters", "method parameters")
       #prop_desc("string", "exceptions", "method exceptions")
       #prop_desc("string", "return", "method return type")
   #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "api"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("method_info")
       #prop_desc("string", "name", "method name")
       #prop_desc("string", "parameters", "method parameters")
       #prop_desc("string", "exceptions", "method exceptions")
       #prop_desc("string", "return", "method return type")
   #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getApiCallList: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
