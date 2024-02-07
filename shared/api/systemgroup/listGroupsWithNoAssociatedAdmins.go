package systemgroup

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of system groups that do not have an administrator.
 (who is not an organization administrator, as they have implicit access to
 system groups) Caller must be an organization administrator.
func ListGroupsWithNoAssociatedAdmins(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
          $ManagedServerGroupSerializer
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "systemgroup"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          $ManagedServerGroupSerializer
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listGroupsWithNoAssociatedAdmins: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
