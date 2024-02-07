package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete the custom values defined for the specified image profile.
 (Note: Attempt to delete values of non-existing keys throws exception. Attempt to
 delete value of existing key which has assigned no values doesn't throw exception.)
func DeleteCustomValues(cnxDetails *api.ConnectionDetails, Label string, $param.getFlagName() $param.getType()) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#return_int_success()](client, "image/profile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteCustomValues: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
