package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Updates repository source URL
func UpdateRepoUrl(cnxDetails *api.ConnectionDetails, Id int, Url string, Label string) (*types.ContentSource, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"id":    Id,
		"url":   Url,
		"label": Label,
	}

	res, err := api.Post[types.ContentSource](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updateRepoUrl: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
