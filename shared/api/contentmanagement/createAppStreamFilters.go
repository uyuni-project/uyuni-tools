package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create Filters for AppStream Modular Channel and attach them to CLM Project
func CreateAppStreamFilters(cnxDetails *api.ConnectionDetails, Prefix string, ChannelLabel string, ProjectLabel string) (*types.#return_array_begin() $ContentFilterSerializer #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"prefix":       Prefix,
		"channelLabel":       ChannelLabel,
		"projectLabel":       ProjectLabel,
	}

	res, err := api.Post[types.#return_array_begin() $ContentFilterSerializer #array_end()](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createAppStreamFilters: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
