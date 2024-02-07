package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Adds an action to remove installed packages on the system to an Action
 Chain.
func AddPackageRemoval(cnxDetails *api.ConnectionDetails, Sid int, PackageIds []int, ChainLabel string) (*types.#param_desc("int", "actionId", "The action id of the scheduled action or exception"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"packageIds":       PackageIds,
		"chainLabel":       ChainLabel,
	}

	res, err := api.Post[types.#param_desc("int", "actionId", "The action id of the scheduled action or exception")](client, "actionchain", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addPackageRemoval: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
