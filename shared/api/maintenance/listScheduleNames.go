package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List Schedule Names visible to user
func ListScheduleNames(cnxDetails *api.ConnectionDetails) (*types.#array_single("string", "maintenance schedule names"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
	}

	res, err := api.Post[types.#array_single("string", "maintenance schedule names")](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listScheduleNames: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
