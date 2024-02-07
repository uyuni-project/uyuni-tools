package external

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get whether we place users into the organization that corresponds
 to the "orgunit" set on the IPA server. The orgunit name must match exactly the
 #product() organization name. Can only be called by a #product() Administrator.
func GetUseOrgUnit(cnxDetails *api.ConnectionDetails) (*types.#param_desc("boolean", "use", "True if we should use the IPA
 orgunit to determine which organization to create the user in, false otherwise"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "user/external"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("boolean", "use", "True if we should use the IPA
 orgunit to determine which organization to create the user in, false otherwise")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getUseOrgUnit: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
