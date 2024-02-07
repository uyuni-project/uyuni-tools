package keys

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// list all keys for the org associated with the user logged into the
             given session
func ListAllKeys(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
          #struct_begin("key")
              #prop("string", "description")
              #prop("string", "type")
          #struct_end()
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/keys"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          #struct_begin("key")
              #prop("string", "description")
              #prop("string", "type")
          #struct_end()
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listAllKeys: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
