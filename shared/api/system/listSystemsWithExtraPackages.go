package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List systems with extra packages
func ListSystemsWithExtraPackages(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
         #struct_begin("system")
             #prop_desc("int", "id", "System ID")
             #prop_desc("string", "name", "System profile name")
             #prop_desc("int", "extra_pkg_count", "Extra packages count")
         #struct_end()
     #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
         #struct_begin("system")
             #prop_desc("int", "id", "System ID")
             #prop_desc("string", "name", "System profile name")
             #prop_desc("int", "extra_pkg_count", "Extra packages count")
         #struct_end()
     #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listSystemsWithExtraPackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
