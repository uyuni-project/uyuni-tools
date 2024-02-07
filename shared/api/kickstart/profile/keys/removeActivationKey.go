package keys

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove an activation key association from the kickstart profile
func RemoveActivationKey(cnxDetails *api.ConnectionDetails, KsLabel string, Key string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
		"key":       Key,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/profile/keys", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeActivationKey: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
