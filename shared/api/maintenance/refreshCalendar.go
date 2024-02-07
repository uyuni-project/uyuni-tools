package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Refresh maintenance calendar data using the configured URL
func RefreshCalendar(cnxDetails *api.ConnectionDetails, Label string, $param.getFlagName() $param.getType()) (*types.#return_array_begin()
       $RescheduleResultSerializer
     #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#return_array_begin()
       $RescheduleResultSerializer
     #array_end()](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute refreshCalendar: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
