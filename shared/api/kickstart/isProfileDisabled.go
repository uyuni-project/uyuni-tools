package kickstart

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns whether a kickstart profile is disabled
func IsProfileDisabled(cnxDetails *api.ConnectionDetails, ProfileLabel string) (*types.#param_desc("boolean", "disabled", "true if profile is disabled"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart"
	params := ""
	if ProfileLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("boolean", "disabled", "true if profile is disabled")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute isProfileDisabled: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
