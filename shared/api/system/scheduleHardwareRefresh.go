package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule a hardware refresh for a system.
func ScheduleHardwareRefresh(cnxDetails *api.ConnectionDetails, Sid int, EarliestOccurrence $date) (*types.#param_desc("int", "actionId", "The action id of the scheduled action"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"earliestOccurrence":       EarliestOccurrence,
	}

	res, err := api.Post[types.#param_desc("int", "actionId", "The action id of the scheduled action")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleHardwareRefresh: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
