package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Updates a ContentSource (repo)
func UpdateRepo(cnxDetails *api.ConnectionDetails, Id int, Label string, Url string) (*types.ContentSource, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"id":       Id,
		"label":       Label,
		"url":       Url,
	}

	res, err := api.Post[types.ContentSource](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updateRepo: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
