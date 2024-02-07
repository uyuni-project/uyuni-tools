package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Adds a filter for a given repo.
func AddRepoFilter(cnxDetails *api.ConnectionDetails, Label string, $param.getFlagName() $param.getType()) (*types.#param_desc("int", "order", "sort order for new filter"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#param_desc("int", "order", "sort order for new filter")](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addRepoFilter: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
