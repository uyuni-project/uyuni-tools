package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the system ID file for a given server.
func DownloadSystemId(cnxDetails *api.ConnectionDetails, Sid int) (*types.#param("string", "id"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
	}

	res, err := api.Post[types.#param("string", "id")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute downloadSystemId: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
