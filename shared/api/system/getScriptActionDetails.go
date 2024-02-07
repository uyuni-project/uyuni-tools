package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns script details for script run actions
func GetScriptActionDetails(cnxDetails *api.ConnectionDetails, ActionId int) (*types.#struct_begin("Script details")
          #prop_desc("int" "id" "action id")
          #prop_desc("string" "content" "script content")
          #prop_desc("string" "run_as_user" "Run as user")
          #prop_desc("string" "run_as_group" "Run as group")
          #prop_desc("int" "timeout" "Timeout in seconds")
          #return_array_begin()
              $ScriptResultSerializer
          #array_end()
      #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if ActionId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("Script details")
          #prop_desc("int" "id" "action id")
          #prop_desc("string" "content" "script content")
          #prop_desc("string" "run_as_user" "Run as user")
          #prop_desc("string" "run_as_group" "Run as group")
          #prop_desc("int" "timeout" "Timeout in seconds")
          #return_array_begin()
              $ScriptResultSerializer
          #array_end()
      #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getScriptActionDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
