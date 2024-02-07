package external

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set whether we place users into the organization that corresponds
 to the "orgunit" set on the IPA server. The orgunit name must match exactly the
 #product() organization name. Can only be called by a #product() Administrator.
func SetUseOrgUnit(cnxDetails *api.ConnectionDetails, UseOrgUnit bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"useOrgUnit":       UseOrgUnit,
	}

	res, err := api.Post[types.#return_int_success()](client, "user/external", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setUseOrgUnit: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
