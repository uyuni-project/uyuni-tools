package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get a list of a kickstart profile's software packages.
func GetSoftwareList(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#array_single("string", "the list of the kickstart profile's software packages"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/profile/software"
	params := ""
	if KsLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("string", "the list of the kickstart profile's software packages")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getSoftwareList: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
