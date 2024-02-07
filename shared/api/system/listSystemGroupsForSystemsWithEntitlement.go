package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns the groups information a system is member of, for all the systems visible to the passed user
 and that are entitled with the passed entitlement.
func ListSystemGroupsForSystemsWithEntitlement(cnxDetails *api.ConnectionDetails, Entitlement string) (*types.#return_array_begin()
     $SystemGroupsDTOSerializer
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if Entitlement {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
     $SystemGroupsDTOSerializer
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listSystemGroupsForSystemsWithEntitlement: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
