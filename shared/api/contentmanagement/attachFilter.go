package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Attach a Filter to a Project
func AttachFilter(cnxDetails *api.ConnectionDetails, ProjectLabel string, FilterId int) (*types.ContentFilter, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"projectLabel": ProjectLabel,
		"filterId":     FilterId,
	}

	res, err := api.Post[types.ContentFilter](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute attachFilter: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
