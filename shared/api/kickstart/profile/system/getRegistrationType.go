package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// returns the registration type of a given kickstart profile.
 Registration Type can be one of reactivation/deletion/none
 These types determine the behaviour of the registration when using
 this profile for reprovisioning.
func GetRegistrationType(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#param("string", "the registration type")
      #options()
         #item ("reactivation")
         #item ("deletion")
         #item ("none")
      #options_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
	}

	res, err := api.Post[types.#param("string", "the registration type")
      #options()
         #item ("reactivation")
         #item ("deletion")
         #item ("none")
      #options_end()](client, "kickstart/profile/system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getRegistrationType: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
