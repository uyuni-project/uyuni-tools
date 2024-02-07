package config

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Given a list of paths and a server, returns details about
 the latest revisions of the paths.
func LookupFileInfo(cnxDetails *api.ConnectionDetails, Sid int, $param.getFlagName() $param.getType()) (*types.#return_array_begin()
          $ConfigRevisionSerializer
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system/config"
	params := ""
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if $param.getFlagName() {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          $ConfigRevisionSerializer
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute lookupFileInfo: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
