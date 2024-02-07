package tree

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the available kickstartable install types (rhel2,3,4,5 and
 fedora9+).
func ListInstallTypes(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin() $KickstartInstallTypeSerializer #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/tree"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin() $KickstartInstallTypeSerializer #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listInstallTypes: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
