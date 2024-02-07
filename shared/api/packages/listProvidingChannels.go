package packages

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the channels that provide the a package.
func ListProvidingChannels(cnxDetails *api.ConnectionDetails, Pid int) (*types.#return_array_begin()
   #struct_begin("channel")
     #prop("string", "label")
     #prop("string", "parent_label")
     #prop("string", "name")
   #struct_end()
 #array_end(), error) {
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

    res, err := api.Get[types.#return_array_begin()
   #struct_begin("channel")
     #prop("string", "label")
     #prop("string", "parent_label")
     #prop("string", "name")
   #struct_end()
 #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listProvidingChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
