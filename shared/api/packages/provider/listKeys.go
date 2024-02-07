package provider

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all security keys associated with a package provider.
 User executing the request must be a #product() administrator.
func ListKeys(cnxDetails *api.ConnectionDetails, ProviderName string) (*types.#return_array_begin()
      $PackageKeySerializer
  #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "packages/provider"
	params := ""
	if ProviderName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
      $PackageKeySerializer
  #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listKeys: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
