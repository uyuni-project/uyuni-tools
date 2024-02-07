package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of empty system profiles visible to user (created by the createSystemProfile method).
func ListEmptySystemProfiles(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
              $EmptySystemProfileSerializer
          #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
              $EmptySystemProfileSerializer
          #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listEmptySystemProfiles: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
