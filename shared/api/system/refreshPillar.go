package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// refresh all the pillar data of a list of systems.
func RefreshPillar(cnxDetails *api.ConnectionDetails, Sids []int, Subset string) (*types.#array_single("int", "skippedIds", "System IDs which couldn't be refreshed"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sids":       Sids,
		"subset":       Subset,
	}

	res, err := api.Post[types.#array_single("int", "skippedIds", "System IDs which couldn't be refreshed")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute refreshPillar: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
