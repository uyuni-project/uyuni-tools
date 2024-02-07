package search

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Search the lucene package indexes for all packages which
          match the given name.
func Name(cnxDetails *api.ConnectionDetails, Name string) (*types.#return_array_begin()
   $PackageOverviewSerializer
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "packages/search"
	params := ""
	if Name {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
   $PackageOverviewSerializer
 #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute name: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
