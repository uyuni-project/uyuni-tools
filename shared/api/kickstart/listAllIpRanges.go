package kickstart

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all Ip Ranges and their associated kickstarts available
 in the user's org.
func ListAllIpRanges(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin() $KickstartIpRangeSerializer #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin() $KickstartIpRangeSerializer #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listAllIpRanges: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
