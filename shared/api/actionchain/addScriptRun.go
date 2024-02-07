package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Add an action with label to run a script to an Action Chain.
 NOTE: The script body must be Base64 encoded!
func AddScriptRun(cnxDetails *api.ConnectionDetails, Sid int, ChainLabel string, ScriptLabel string, Uid string, Gid string, Timeout int, ScriptBody string) (*types.#param_desc("int", "actionId", "The id of the action or throw an exception"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"chainLabel":       ChainLabel,
		"scriptLabel":       ScriptLabel,
		"uid":       Uid,
		"gid":       Gid,
		"timeout":       Timeout,
		"scriptBody":       ScriptBody,
	}

	res, err := api.Post[types.#param_desc("int", "actionId", "The id of the action or throw an exception")](client, "actionchain", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addScriptRun: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
