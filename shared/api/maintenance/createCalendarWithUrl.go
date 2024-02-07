package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new maintenance calendar
func CreateCalendarWithUrl(cnxDetails *api.ConnectionDetails, Label string, Url string) (*types.#return_array_begin()
 $MaintenanceCalendarSerializer
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"url":       Url,
	}

	res, err := api.Post[types.#return_array_begin()
 $MaintenanceCalendarSerializer
 #array_end()](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createCalendarWithUrl: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
