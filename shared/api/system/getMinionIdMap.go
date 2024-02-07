package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Return a map from Salt minion IDs to System IDs.
 Map entries are limited to systems that are visible by the current user.
func GetMinionIdMap(cnxDetails *api.ConnectionDetails) (*types.#param_desc("map", "id_map", "minion IDs to system IDs"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("map", "id_map", "minion IDs to system IDs")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getMinionIdMap: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
