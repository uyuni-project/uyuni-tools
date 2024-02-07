package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove a Content Filter
func RemoveFilter(cnxDetails *api.ConnectionDetails, FilterId int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"filterId":       FilterId,
	}

	res, err := api.Post[types.#return_int_success()](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeFilter: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
