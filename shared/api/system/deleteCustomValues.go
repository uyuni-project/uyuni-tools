package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete the custom values defined for the custom system information keys
 provided from the given system.
 (Note: Attempt to delete values of non-existing keys throws exception. Attempt to
 delete value of existing key which has assigned no values doesn't throw exception.)
func DeleteCustomValues(cnxDetails *api.ConnectionDetails, Sid int, Keys []string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"keys":       Keys,
	}

	res, err := api.Post[types.#return_int_success()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteCustomValues: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
