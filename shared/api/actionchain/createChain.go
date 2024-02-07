package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create an Action Chain.
func CreateChain(cnxDetails *api.ConnectionDetails, ChainLabel string) (*types.#param_desc("int", "actionId", "The ID of the created action chain"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"chainLabel":       ChainLabel,
	}

	res, err := api.Post[types.#param_desc("int", "actionId", "The ID of the created action chain")](client, "actionchain", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createChain: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
