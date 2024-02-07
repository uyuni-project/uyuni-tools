package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Fetch results from a script execution. Returns an empty array if no
 results are yet available.
func GetScriptResults(cnxDetails *api.ConnectionDetails, ActionId int) (*types.#return_array_begin()
              $ScriptResultSerializer
         #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"actionId":       ActionId,
	}

	res, err := api.Post[types.#return_array_begin()
              $ScriptResultSerializer
         #array_end()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getScriptResults: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
