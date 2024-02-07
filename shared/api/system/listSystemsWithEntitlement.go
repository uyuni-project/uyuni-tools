package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lists the systems that have the given entitlement
func ListSystemsWithEntitlement(cnxDetails *api.ConnectionDetails, EntitlementName string) (*types.#return_array_begin()
                  $SystemOverviewSerializer
              #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if EntitlementName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
                  $SystemOverviewSerializer
              #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listSystemsWithEntitlement: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
