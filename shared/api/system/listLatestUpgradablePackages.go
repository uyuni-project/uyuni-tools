package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the list of latest upgradable packages for a given system.
func ListLatestUpgradablePackages(cnxDetails *api.ConnectionDetails, Sid int) (*types.#return_array_begin()
      #struct_begin("package")
          #prop("string", "name")
          #prop("string", "arch")
          #prop("string", "from_version")
          #prop("string", "from_release")
          #prop("string", "from_epoch")
          #prop("string", "to_version")
          #prop("string", "to_release")
          #prop("string", "to_epoch")
          #prop("string", "to_package_id")
      #struct_end()
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
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
      #struct_begin("package")
          #prop("string", "name")
          #prop("string", "arch")
          #prop("string", "from_version")
          #prop("string", "from_release")
          #prop("string", "from_epoch")
          #prop("string", "to_version")
          #prop("string", "to_release")
          #prop("string", "to_epoch")
          #prop("string", "to_package_id")
      #struct_end()
 #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listLatestUpgradablePackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
