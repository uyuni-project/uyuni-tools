package auth

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Login using a username and password. Returns the session key
 used by most other API methods.
func Login(cnxDetails *api.ConnectionDetails, Username string, Password string, Duration int) (*types.#session_key(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"username":       Username,
		"password":       Password,
		"duration":       Duration,
	}

	res, err := api.Post[types.#session_key()](client, "auth", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute login: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
