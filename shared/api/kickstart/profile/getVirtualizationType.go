package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// For given kickstart profile label returns label of
 virtualization type it's using
func GetVirtualizationType(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#param_desc("string", "virtLabel",
 "Label of virtualization type."), error) {
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

    res, err := api.Get[types.#param_desc("string", "virtLabel",
 "Label of virtualization type.")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getVirtualizationType: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
