package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lookup a specific maintenance schedule
func GetScheduleDetails(cnxDetails *api.ConnectionDetails, Name string) (*types.#return_array_begin()
 $MaintenanceScheduleSerializer
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
	}

	res, err := api.Post[types.#return_array_begin()
 $MaintenanceScheduleSerializer
 #array_end()](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getScheduleDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
