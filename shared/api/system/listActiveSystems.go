package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of active servers visible to the user.
func ListActiveSystems(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
          $ShortSystemInfoSerializer
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
          $ShortSystemInfoSerializer
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listActiveSystems: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
