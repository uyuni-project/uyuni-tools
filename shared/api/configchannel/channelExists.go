package configchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Check for the existence of the config channel provided.
func ChannelExists(cnxDetails *api.ConnectionDetails, Label string) (*types.#param_desc("int", "existence", "1 if exists, 0 otherwise"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
	}

	res, err := api.Post[types.#param_desc("int", "existence", "1 if exists, 0 otherwise")](client, "configchannel", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute channelExists: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
