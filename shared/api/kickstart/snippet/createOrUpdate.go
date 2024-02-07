package snippet

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Will create a snippet with the given name and contents if it
      doesn't exist. If it does exist, the existing snippet will be updated.
func CreateOrUpdate(cnxDetails *api.ConnectionDetails, Name string, Contents string) (*types.Snippet, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
		"contents":       Contents,
	}

	res, err := api.Post[types.Snippet](client, "kickstart/snippet", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createOrUpdate: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
