package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of all errata that are relevant to the system.
func GetRelevantErrata(cnxDetails *api.ConnectionDetails, Sid int, Sids []int) (*types.#return_array_begin()
          $ErrataOverviewSerializer
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if Sids {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          $ErrataOverviewSerializer
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getRelevantErrata: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
