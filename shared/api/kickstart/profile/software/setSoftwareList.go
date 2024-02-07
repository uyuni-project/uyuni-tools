package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set the list of software packages for a kickstart profile.
func SetSoftwareList(cnxDetails *api.ConnectionDetails, KsLabel string, $param.getFlagName() $param.getType(), $param.getFlagName() $param.getType(), IgnoreMissing bool, NoBase bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
		"$param.getName()":       $param.getFlagName(),
		"$param.getName()":       $param.getFlagName(),
		"ignoreMissing":       IgnoreMissing,
		"noBase":       NoBase,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/profile/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setSoftwareList: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
