package tree

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete a Kickstart Tree (Distribution) from #product().
func Delete(cnxDetails *api.ConnectionDetails, TreeLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"treeLabel":       TreeLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/tree", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute delete: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
