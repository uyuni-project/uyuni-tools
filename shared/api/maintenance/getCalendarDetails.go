package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lookup a specific maintenance schedule
func GetCalendarDetails(cnxDetails *api.ConnectionDetails, Label string) (*types.#return_array_begin()
 $MaintenanceCalendarSerializer
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
	}

	res, err := api.Post[types.#return_array_begin()
 $MaintenanceCalendarSerializer
 #array_end()](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getCalendarDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
