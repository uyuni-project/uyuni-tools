package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Provides a list of packages installed on a system that are also
          contained in the given channel.  The installed package list did not
          include arch information before RHEL 5, so it is arch unaware.  RHEL 5
          systems do upload the arch information, and thus are arch aware.
func ListPackagesFromChannel(cnxDetails *api.ConnectionDetails, Sid int, ChannelLabel string) (*types.#return_array_begin()
      $PackageSerializer
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
	if ChannelLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
      $PackageSerializer
  #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listPackagesFromChannel: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
