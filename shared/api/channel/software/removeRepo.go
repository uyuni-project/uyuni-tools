package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Removes a repository
func RemoveRepo(cnxDetails *api.ConnectionDetails, Id long, Label string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"id":       Id,
		"label":       Label,
	}

	res, err := api.Post[types.#return_int_success()](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeRepo: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
