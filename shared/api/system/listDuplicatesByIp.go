package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List duplicate systems by IP Address.
func ListDuplicatesByIp(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
           #struct_begin("Duplicate Group")
                   #prop("string", "ip")
                   #prop_array_begin("systems")
                      $NetworkDtoSerializer
                   #prop_array_end()
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
           #struct_begin("Duplicate Group")
                   #prop("string", "ip")
                   #prop_array_begin("systems")
                      $NetworkDtoSerializer
                   #prop_array_end()
           #struct_end()
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listDuplicatesByIp: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
