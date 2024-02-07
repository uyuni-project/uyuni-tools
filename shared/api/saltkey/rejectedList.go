package saltkey

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List of rejected salt keys
func RejectedList(cnxDetails *api.ConnectionDetails) (*types.#array_single("string", "Rejected salt key list"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "saltkey"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("string", "Rejected salt key list")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rejectedList: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
