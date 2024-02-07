package powermanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Execute power management action 'Reboot'
func Reboot(cnxDetails *api.ConnectionDetails, Sid int, Name string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"name":       Name,
	}

	res, err := api.Post[types.#return_int_success()](client, "system/provisioning/powermanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute reboot: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
