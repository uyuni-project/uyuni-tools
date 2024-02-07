package api

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lists all available api calls for the specified namespace
func GetApiNamespaceCallList(cnxDetails *api.ConnectionDetails, Namespace string) (*types.#struct_begin("method_info")
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
	if Namespace {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
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
		return nil, fmt.Errorf("failed to execute getApiNamespaceCallList: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
