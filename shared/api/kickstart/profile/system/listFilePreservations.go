package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns the set of all file preservations associated with the given
 kickstart profile.
func ListFilePreservations(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#return_array_begin()
         $FileListSerializer
     #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/profile/system"
	params := ""
	if KsLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
         $FileListSerializer
     #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listFilePreservations: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
