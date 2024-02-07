package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new maintenance Schedule
func CreateSchedule(cnxDetails *api.ConnectionDetails, Name string, Type string, Calendar string) (*types.#return_array_begin()
 $MaintenanceScheduleSerializer
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
		"type":       Type,
		"calendar":       Calendar,
	}

	res, err := api.Post[types.#return_array_begin()
 $MaintenanceScheduleSerializer
 #array_end()](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createSchedule: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
