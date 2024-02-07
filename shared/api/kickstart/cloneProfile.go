package kickstart

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Clone a Kickstart Profile
func CloneProfile(cnxDetails *api.ConnectionDetails, KsLabelToClone string, NewKsLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabelToClone":       KsLabelToClone,
		"newKsLabel":       NewKsLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute cloneProfile: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
