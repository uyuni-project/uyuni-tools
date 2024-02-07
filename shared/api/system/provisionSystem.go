package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Provision a system using the specified kickstart/autoinstallation profile.
func ProvisionSystem(cnxDetails *api.ConnectionDetails, Sid int, ProfileName string, EarliestDate $date) (*types.#param_desc("int", "id", "ID of the action scheduled, otherwise exception thrown
 on error"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"profileName":       ProfileName,
		"earliestDate":       EarliestDate,
	}

	res, err := api.Post[types.#param_desc("int", "id", "ID of the action scheduled, otherwise exception thrown
 on error")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute provisionSystem: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
