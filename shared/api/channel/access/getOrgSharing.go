package access

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get organization sharing access control.
func GetOrgSharing(cnxDetails *api.ConnectionDetails, ChannelLabel string) (*types.#param_desc("string", "access", "The access value (one of the following: 'public', 'private', or 'protected'"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "channel/access"
	params := ""
	if ChannelLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("string", "access", "The access value (one of the following: 'public', 'private', or 'protected'")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getOrgSharing: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
