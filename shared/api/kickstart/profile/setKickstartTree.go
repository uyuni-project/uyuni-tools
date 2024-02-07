package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set the kickstart tree for a kickstart profile.
func SetKickstartTree(cnxDetails *api.ConnectionDetails, KsLabel string, KstreeLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
		"kstreeLabel":       KstreeLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/profile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setKickstartTree: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
