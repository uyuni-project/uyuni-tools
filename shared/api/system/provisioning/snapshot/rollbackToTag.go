package snapshot

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Rollbacks server to snapshot
func RollbackToTag(cnxDetails *api.ConnectionDetails, Sid int, TagName string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"tagName":       TagName,
	}

	res, err := api.Post[types.#return_int_success()](client, "system/provisioning/snapshot", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rollbackToTag: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
