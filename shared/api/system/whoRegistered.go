package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns information about the user who registered the system
func WhoRegistered(cnxDetails *api.ConnectionDetails, Sid int) (*types.User, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
	}

	res, err := api.Post[types.User](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute whoRegistered: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
