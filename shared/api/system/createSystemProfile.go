package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Creates a system record in database for a system that is not registered.
 Either "hwAddress" or "hostname" prop must be specified in the "data" struct.
 If a system(s) matching given data exists, a SystemsExistFaultException is thrown which
 contains matching system IDs in its message.
func CreateSystemProfile(cnxDetails *api.ConnectionDetails, SystemName string, $param.getFlagName() $param.getType()) (*types.#param_desc("int", "systemId", "The id of the created system"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"systemName":       SystemName,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#param_desc("int", "systemId", "The id of the created system")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createSystemProfile: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
