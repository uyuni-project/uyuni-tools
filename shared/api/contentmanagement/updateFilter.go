package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update a Content Filter
 #paragraph_end()
 #paragraph()
 See also: createFilter(), listFilterCriteria()
func UpdateFilter(cnxDetails *api.ConnectionDetails, FilterId int, Name string, Rule string, $param.getFlagName() $param.getType()) (*types.ContentFilter, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"filterId":       FilterId,
		"name":       Name,
		"rule":       Rule,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.ContentFilter](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updateFilter: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
