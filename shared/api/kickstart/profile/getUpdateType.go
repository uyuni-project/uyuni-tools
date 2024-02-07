package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the update type for a kickstart profile.
func GetUpdateType(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#param_desc("string", "update_type", "Update type for this Kickstart Profile."), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/profile"
	params := ""
	if KsLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("string", "update_type", "Update type for this Kickstart Profile.")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getUpdateType: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
