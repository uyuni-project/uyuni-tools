package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new stored Package Profile from a systems
      installed package list.
func CreatePackageProfile(cnxDetails *api.ConnectionDetails, Sid int, ProfileLabel string, Description string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"profileLabel":       ProfileLabel,
		"description":       Description,
	}

	res, err := api.Post[types.#return_int_success()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createPackageProfile: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
