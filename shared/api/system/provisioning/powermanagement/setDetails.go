package powermanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get current power management settings of the given system
func SetDetails(cnxDetails *api.ConnectionDetails, Sid int, $param.getFlagName() $param.getType(), Name string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"$param.getName()":       $param.getFlagName(),
		"name":       Name,
	}

	res, err := api.Post[types.#return_int_success()](client, "system/provisioning/powermanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
