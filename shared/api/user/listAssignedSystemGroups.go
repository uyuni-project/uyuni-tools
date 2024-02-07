package user

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns the system groups that a user can administer.
func ListAssignedSystemGroups(cnxDetails *api.ConnectionDetails, Login string) (*types.#return_array_begin()
     #struct_begin("system group")
       #prop("int", "id")
       #prop("string", "name")
       #prop("string", "description")
       #prop("int", "system_count")
       #prop_desc("int", "org_id", "Organization ID for this system group.")
     #struct_end()
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "user"
	params := ""
	if Login {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
     #struct_begin("system group")
       #prop("int", "id")
       #prop("string", "name")
       #prop("string", "description")
       #prop("int", "system_count")
       #prop_desc("int", "org_id", "Organization ID for this system group.")
     #struct_end()
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listAssignedSystemGroups: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
