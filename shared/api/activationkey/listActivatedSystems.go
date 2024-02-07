package activationkey

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the systems activated with the key provided.
func ListActivatedSystems(cnxDetails *api.ConnectionDetails, Key string) (*types.#return_array_begin()
       #struct_begin("system structure")
           #prop_desc("int", "id", "System id")
           #prop("string", "hostname")
           #prop_desc("$date",  "last_checkin", "Last time server
               successfully checked in")
       #struct_end()
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "activationkey"
	params := ""
	if Key {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
       #struct_begin("system structure")
           #prop_desc("int", "id", "System id")
           #prop("string", "hostname")
           #prop_desc("$date",  "last_checkin", "Last time server
               successfully checked in")
       #struct_end()
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listActivatedSystems: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
