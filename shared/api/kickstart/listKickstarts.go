package kickstart

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Provides a list of kickstart profiles visible to the user's
 org
func ListKickstarts(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin() $KickstartDtoSerializer #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin() $KickstartDtoSerializer #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listKickstarts: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
