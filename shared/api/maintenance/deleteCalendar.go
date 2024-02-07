package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove a maintenance calendar
func DeleteCalendar(cnxDetails *api.ConnectionDetails, Label string, CancelScheduledActions bool) (*types.#return_array_begin()
       $RescheduleResultSerializer
     #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"cancelScheduledActions":       CancelScheduledActions,
	}

	res, err := api.Post[types.#return_array_begin()
       $RescheduleResultSerializer
     #array_end()](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteCalendar: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
