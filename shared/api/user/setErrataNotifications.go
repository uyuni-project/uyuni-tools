package user

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Enables/disables errata mail notifications for a specific user.
func SetErrataNotifications(cnxDetails *api.ConnectionDetails, Login string, Value bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"login":       Login,
		"value":       Value,
	}

	res, err := api.Post[types.#return_int_success()](client, "user", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setErrataNotifications: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
