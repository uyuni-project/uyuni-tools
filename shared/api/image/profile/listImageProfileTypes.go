package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List available image store types
func ListImageProfileTypes(cnxDetails *api.ConnectionDetails) (*types.#array_single("string", "the list of image profile types"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "image/profile"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("string", "the list of image profile types")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listImageProfileTypes: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
