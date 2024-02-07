package content

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List organization credentials (mirror credentials) available in
             #product().
func ListCredentials(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
                       $MirrorCredentialsDtoSerializer
                    #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "sync/content"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
                       $MirrorCredentialsDtoSerializer
                    #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listCredentials: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
