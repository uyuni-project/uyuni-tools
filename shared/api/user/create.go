package user

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new user.
func Create(cnxDetails *api.ConnectionDetails, Login string, Password string, FirstName string, LastName string, Email string, Login string, UsePamAuth int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"login":       Login,
		"password":       Password,
		"firstName":       FirstName,
		"lastName":       LastName,
		"email":       Email,
		"login":       Login,
		"usePamAuth":       UsePamAuth,
	}

	res, err := api.Post[types.#return_int_success()](client, "user", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
