package config

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new file (text or binary) or directory with the given path, or
 update an existing path on a server.
func CreateOrUpdatePath(cnxDetails *api.ConnectionDetails, Sid int, Path string) (*types.ConfigRevision, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"path":       Path,
	}

	res, err := api.Post[types.ConfigRevision](client, "system/config", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createOrUpdatePath: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
