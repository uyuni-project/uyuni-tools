package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Refresh pillar data and then schedule channels state on the provided systems
func ApplyChannelState(cnxDetails *api.ConnectionDetails, Sids []int) (*types.#array_single("int", "actionId"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sids":       Sids,
	}

	res, err := api.Post[types.#array_single("int", "actionId")](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute applyChannelState: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
