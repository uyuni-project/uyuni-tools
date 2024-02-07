package external

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set whether we should keeps roles assigned to users because of
 their IPA groups even after they log in through a non-IPA method. Can only be
 called by a #product() Administrator.
func SetKeepTemporaryRoles(cnxDetails *api.ConnectionDetails, KeepRoles bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"keepRoles":       KeepRoles,
	}

	res, err := api.Post[types.#return_int_success()](client, "user/external", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setKeepTemporaryRoles: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
