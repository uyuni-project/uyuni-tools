package content

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete organization credentials (mirror credentials) from #product().
func DeleteCredentials(cnxDetails *api.ConnectionDetails, Username string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"username":       Username,
	}

	res, err := api.Post[types.#return_int_success()](client, "sync/content", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteCredentials: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
