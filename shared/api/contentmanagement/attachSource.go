package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Attach a Source to a Project
func AttachSource(cnxDetails *api.ConnectionDetails, ProjectLabel string, SourceType string, SourceLabel string, SourcePosition int) (*types.ContentProjectSource, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"projectLabel":   ProjectLabel,
		"sourceType":     SourceType,
		"sourceLabel":    SourceLabel,
		"sourcePosition": SourcePosition,
	}

	res, err := api.Post[types.ContentProjectSource](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute attachSource: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
