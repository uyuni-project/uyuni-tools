package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the available groups for a given system.
func ListGroups(cnxDetails *api.ConnectionDetails, Sid int) (*types.#return_array_begin()
      #struct_begin("system group")
          #prop_desc("int", "id", "server group id")
          #prop_desc("int", "subscribed", "1 if the given server is subscribed
               to this server group, 0 otherwise")
          #prop_desc("string", "system_group_name", "Name of the server group")
          #prop_desc("string", "sgid", "server group id (Deprecated)")
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
      #struct_begin("system group")
          #prop_desc("int", "id", "server group id")
          #prop_desc("int", "subscribed", "1 if the given server is subscribed
               to this server group, 0 otherwise")
          #prop_desc("string", "system_group_name", "Name of the server group")
          #prop_desc("string", "sgid", "server group id (Deprecated)")
      #struct_end()
  #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listGroups: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
