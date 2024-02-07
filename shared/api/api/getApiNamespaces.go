package api

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lists available API namespaces
func GetApiNamespaces(cnxDetails *api.ConnectionDetails) (*types.#struct_begin("namespace")
        #prop_desc("string", "namespace", "API namespace")
        #prop_desc("string", "handler", "API Handler")
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

    res, err := api.Get[types.#struct_begin("namespace")
        #prop_desc("string", "namespace", "API namespace")
        #prop_desc("string", "handler", "API Handler")
   #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getApiNamespaces: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
