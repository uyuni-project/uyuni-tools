package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the installed packages for a given system.
func ListInstalledPackages(cnxDetails *api.ConnectionDetails, Sid int) (*types.#return_array_begin()
          #struct_begin("package")
                 #prop_desc("int", "package_id", "PackageID, -1 if package is installed but not available in
                 subscribed channels")
                 #prop("string", "name")
                 #prop("string", "epoch")
                 #prop("string", "version")
                 #prop("string", "release")
                 #prop_desc("string", "arch", "architecture label")
                 #prop_desc("date", "installtime", "returned only if known")
                 #prop("boolean", "retracted")
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
                 #prop_desc("int", "package_id", "PackageID, -1 if package is installed but not available in
                 subscribed channels")
                 #prop("string", "name")
                 #prop("string", "epoch")
                 #prop("string", "version")
                 #prop("string", "release")
                 #prop_desc("string", "arch", "architecture label")
                 #prop_desc("date", "installtime", "returned only if known")
                 #prop("boolean", "retracted")
          #struct_end()
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listInstalledPackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
