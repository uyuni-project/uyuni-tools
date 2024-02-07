package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Detach a Source from a Project
func DetachSource(cnxDetails *api.ConnectionDetails, ProjectLabel string, SourceType string, SourceLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"projectLabel":       ProjectLabel,
		"sourceType":       SourceType,
		"sourceLabel":       SourceLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute detachSource: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
