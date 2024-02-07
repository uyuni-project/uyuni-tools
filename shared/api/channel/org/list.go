package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the organizations associated with the given channel
 that may be trusted.
func List(cnxDetails *api.ConnectionDetails, Label string) (*types.#return_array_begin()
      #struct_begin("org")
          #prop("int", "org_id")
          #prop("string", "org_name")
          #prop("boolean", "access_enabled")
     #struct_end()
  #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
	}

	res, err := api.Post[types.#return_array_begin()
      #struct_begin("org")
          #prop("int", "org_id")
          #prop("string", "org_name")
          #prop("boolean", "access_enabled")
     #struct_end()
  #array_end()](client, "channel/org", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute list: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
