package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List schedule names visible to user
func ListCalendarLabels(cnxDetails *api.ConnectionDetails) (*types.#array_single("string", "maintenance calendar labels"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
	}

	res, err := api.Post[types.#array_single("string", "maintenance calendar labels")](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listCalendarLabels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
