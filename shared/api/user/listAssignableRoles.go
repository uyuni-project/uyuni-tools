package user

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of user roles that this user can assign to others.
func ListAssignableRoles(cnxDetails *api.ConnectionDetails) (*types.#array_single("string", "(role label)"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "user"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("string", "(role label)")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listAssignableRoles: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
