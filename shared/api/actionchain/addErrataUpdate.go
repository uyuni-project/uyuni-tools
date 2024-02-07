package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Adds Errata update to an Action Chain.
func AddErrataUpdate(cnxDetails *api.ConnectionDetails, Sid int, $param.getFlagName() $param.getType(), ChainLabel string, $param.getFlagName() $param.getType(), OnlyRelevant bool) (*types.#param_desc("int", "actionId", "The action id of the scheduled action"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"$param.getName()":       $param.getFlagName(),
		"chainLabel":       ChainLabel,
		"$param.getName()":       $param.getFlagName(),
		"onlyRelevant":       OnlyRelevant,
	}

	res, err := api.Post[types.#param_desc("int", "actionId", "The action id of the scheduled action")](client, "actionchain", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addErrataUpdate: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
