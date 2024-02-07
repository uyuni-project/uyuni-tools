package snapshot

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Deletes a snapshot with the given snapshot id
func DeleteSnapshot(cnxDetails *api.ConnectionDetails, SnapId int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"snapId":       SnapId,
	}

	res, err := api.Post[types.#return_int_success()](client, "system/provisioning/snapshot", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteSnapshot: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
