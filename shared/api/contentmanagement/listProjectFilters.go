package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all Filters associated with a Project
func ListProjectFilters(cnxDetails *api.ConnectionDetails, ProjectLabel string) (*types.#return_array_begin()
 $ContentProjectFilterSerializer
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "contentmanagement"
	params := ""
	if ProjectLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
 $ContentProjectFilterSerializer
 #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listProjectFilters: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
