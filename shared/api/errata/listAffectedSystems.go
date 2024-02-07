package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Return the list of systems affected by the errata with the given advisory name.
 For those errata that are present in both vendor and user organizations under the same advisory name,
 this method retrieves the affected systems by both of them.
func ListAffectedSystems(cnxDetails *api.ConnectionDetails, AdvisoryName string) (*types.#return_array_begin()
          $SystemOverviewSerializer
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "errata"
	params := ""
	if AdvisoryName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          $SystemOverviewSerializer
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listAffectedSystems: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
