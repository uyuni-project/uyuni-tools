package external

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get a representation of the role mapping for an external group.
 Can only be called by a #product() Administrator.
func GetExternalGroupToRoleMap(cnxDetails *api.ConnectionDetails, Name string) (*types.UserExtGroup, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "user/external"
	params := ""
	if Name {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.UserExtGroup](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getExternalGroupToRoleMap: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
