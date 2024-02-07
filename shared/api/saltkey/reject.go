package saltkey

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Reject a minion key
func Reject(cnxDetails *api.ConnectionDetails, MinionId string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"minionId":       MinionId,
	}

	res, err := api.Post[types.#return_int_success()](client, "saltkey", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute reject: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
