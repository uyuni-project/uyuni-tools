package external

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get whether we should keeps roles assigned to users because of
 their IPA groups even after they log in through a non-IPA method. Can only be
 called by a #product() Administrator.
func GetKeepTemporaryRoles(cnxDetails *api.ConnectionDetails) (*types.#param_desc("boolean", "keep", "True if we should keep roles
 after users log in through non-IPA method, false otherwise"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "user/external"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("boolean", "keep", "True if we should keep roles
 after users log in through non-IPA method, false otherwise")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getKeepTemporaryRoles: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
