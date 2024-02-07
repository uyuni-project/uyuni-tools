package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update a maintenance schedule
func UpdateSchedule(cnxDetails *api.ConnectionDetails, Name string, $param.getFlagName() $param.getType(), $param.getFlagName() $param.getType()) (*types.RescheduleResult, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
		"$param.getName()":       $param.getFlagName(),
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.RescheduleResult](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updateSchedule: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
