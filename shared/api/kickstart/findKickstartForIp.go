package kickstart

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Find an associated kickstart for a given ip address.
func FindKickstartForIp(cnxDetails *api.ConnectionDetails, IpAddress string) (*types.#param_desc("string", "label", "label of the kickstart. Empty string if not found"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart"
	params := ""
	if IpAddress {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("string", "label", "label of the kickstart. Empty string if not found")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute findKickstartForIp: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
