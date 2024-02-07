package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a Content Environment and appends it behind given Content Environment
func CreateEnvironment(cnxDetails *api.ConnectionDetails, ProjectLabel string, PredecessorLabel string, EnvLabel string, Name string, Description string) (*types.ContentEnvironment, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"projectLabel":       ProjectLabel,
		"predecessorLabel":       PredecessorLabel,
		"envLabel":       EnvLabel,
		"name":       Name,
		"description":       Description,
	}

	res, err := api.Post[types.ContentEnvironment](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createEnvironment: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
