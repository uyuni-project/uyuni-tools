package custom

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all the custom states available to the user.
func ListAvailable(cnxDetails *api.ConnectionDetails) (*types.#array_single("string", "the list of custom channels available to the user"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "recurring/custom"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("string", "the list of custom channels available to the user")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listAvailable: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
