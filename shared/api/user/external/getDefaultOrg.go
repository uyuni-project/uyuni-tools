package external

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the default org that users should be added in if orgunit from
 IPA server isn't found or is disabled. Can only be called by a #product() Administrator.
func GetDefaultOrg(cnxDetails *api.ConnectionDetails) (*types.#param_desc("int", "id", "ID of the default organization. 0 if there is no default"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "user/external"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("int", "id", "ID of the default organization. 0 if there is no default")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getDefaultOrg: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
