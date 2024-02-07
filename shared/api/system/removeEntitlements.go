package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove addon entitlements from a server. Entitlements a server does
 not have are quietly ignored.
func RemoveEntitlements(cnxDetails *api.ConnectionDetails, Sid int, Entitlements []string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"entitlements":       Entitlements,
	}

	res, err := api.Post[types.#return_int_success()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeEntitlements: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
