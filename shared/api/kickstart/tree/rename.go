package tree

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Rename a Kickstart Tree (Distribution) in #product().
func Rename(cnxDetails *api.ConnectionDetails, OriginalLabel string, NewLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"originalLabel":       OriginalLabel,
		"newLabel":       NewLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/tree", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rename: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
