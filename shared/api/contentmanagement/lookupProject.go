package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Look up Content Project with given label
func LookupProject(cnxDetails *api.ConnectionDetails, ProjectLabel string) (*types.ContentProject, error) {
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

    res, err := api.Get[types.ContentProject](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute lookupProject: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
