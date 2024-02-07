package packages

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Retrieve the url that can be used to download a package.
      This will expire after a certain time period.
func GetPackageUrl(cnxDetails *api.ConnectionDetails, Pid int) (*types.string - the download url, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "packages"
	params := ""
	if Pid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.string - the download url](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getPackageUrl: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
