package trusts

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lists all software channels that the organization given is providing to
 the user's organization.
func ListChannelsProvided(cnxDetails *api.ConnectionDetails, OrgId int) (*types.#return_array_begin()
         #struct_begin("channel info")
             #prop("int", "channel_id")
             #prop("string", "channel_name")
             #prop("int", "packages")
             #prop("int", "systems")
         #struct_end()
     #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "org/trusts"
	params := ""
	if OrgId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
         #struct_begin("channel info")
             #prop("int", "channel_id")
             #prop("string", "channel_name")
             #prop("int", "packages")
             #prop("int", "systems")
         #struct_end()
     #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listChannelsProvided: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
