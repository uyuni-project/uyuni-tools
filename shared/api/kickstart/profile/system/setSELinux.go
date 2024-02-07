package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Sets the SELinux enforcing mode property of a kickstart profile
 so that a system created using this profile will be have
 the appropriate SELinux enforcing mode.
func SetSELinux(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/profile/system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setSELinux: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
