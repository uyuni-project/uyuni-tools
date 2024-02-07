package tree

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete a kickstarttree and any profiles associated with
 this kickstart tree.  WARNING:  This will delete all profiles
 associated with this kickstart tree!
func DeleteTreeAndProfiles(cnxDetails *api.ConnectionDetails, TreeLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"treeLabel":       TreeLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/tree", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteTreeAndProfiles: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
