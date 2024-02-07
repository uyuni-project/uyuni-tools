package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List of available filter criteria
func ListFilterCriteria(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
 #struct_begin("Filter Criteria")
 #prop("string", "type")
 #prop("string", "matcher")
 #prop("string", "field")
 #struct_end()
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "contentmanagement"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
 #struct_begin("Filter Criteria")
 #prop("string", "type")
 #prop("string", "matcher")
 #prop("string", "field")
 #struct_end()
 #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listFilterCriteria: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
