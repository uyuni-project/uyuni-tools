package packages

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all source packages in user's organization.
func ListSourcePackages(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
   #struct_begin("source_package")
     #prop("int", "id")
     #prop("string", "name")
   #struct_end()
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "packages"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
   #struct_begin("source_package")
     #prop("int", "id")
     #prop("string", "name")
   #struct_end()
 #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listSourcePackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
