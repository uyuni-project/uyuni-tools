package systemgroup

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns the list of users who can administer the given group.
 Caller must be a system group admin or an organization administrator.
func ListAdministrators(cnxDetails *api.ConnectionDetails, SystemGroupName string) (*types.#return_array_begin()
      $UserSerializer
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "systemgroup"
	params := ""
	if SystemGroupName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
      $UserSerializer
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listAdministrators: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
