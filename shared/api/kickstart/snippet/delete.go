package snippet

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete the specified snippet.
      If the snippet is not found, 0 is returned.
func Delete(cnxDetails *api.ConnectionDetails, Name string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/snippet", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute delete: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
