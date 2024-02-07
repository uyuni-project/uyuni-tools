package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Replaces the existing set of filters for a given repo.
 Filters are ranked by their order in the array.
func SetRepoFilters(cnxDetails *api.ConnectionDetails, Label string, $param.getFlagName() $param.getType()) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#return_int_success()](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setRepoFilters: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
