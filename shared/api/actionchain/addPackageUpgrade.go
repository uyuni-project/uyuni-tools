package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Adds an action to upgrade installed packages on the system to an Action
 Chain.
func AddPackageUpgrade(cnxDetails *api.ConnectionDetails, Sid int, PackageIds []int, ChainLabel string) (*types.#param_desc("int", "actionId", "The id of the action or throw an exception"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"packageIds":       PackageIds,
		"chainLabel":       ChainLabel,
	}

	res, err := api.Post[types.#param_desc("int", "actionId", "The id of the action or throw an exception")](client, "actionchain", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addPackageUpgrade: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
