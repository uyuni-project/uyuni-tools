package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Adds an action to verify installed packages on the system to an Action
 Chain.
func AddPackageVerify(cnxDetails *api.ConnectionDetails, Sid int, PackageIds []int, ChainLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"packageIds":       PackageIds,
		"chainLabel":       ChainLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "actionchain", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addPackageVerify: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
