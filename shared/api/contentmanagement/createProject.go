package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create Content Project
func CreateProject(cnxDetails *api.ConnectionDetails, ProjectLabel string, Name string, Description string) (*types.ContentProject, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"projectLabel": ProjectLabel,
		"name":         Name,
		"description":  Description,
	}

	res, err := api.Post[types.ContentProject](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createProject: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
