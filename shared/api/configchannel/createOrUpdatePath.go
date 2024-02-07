package configchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new file or directory with the given path, or
 update an existing path.
func CreateOrUpdatePath(cnxDetails *api.ConnectionDetails, Label string, Path string, IsDir bool) (*types.ConfigRevision, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"path":       Path,
		"isDir":       IsDir,
	}

	res, err := api.Post[types.ConfigRevision](client, "configchannel", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createOrUpdatePath: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
