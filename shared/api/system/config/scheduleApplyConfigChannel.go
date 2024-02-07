package config

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule highstate application for a given system.
func ScheduleApplyConfigChannel(cnxDetails *api.ConnectionDetails, Sids []int, EarliestOccurrence $date, Test bool) (*types.#param("int", "actionId", "The action id of the scheduled action"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sids":       Sids,
		"earliestOccurrence":       EarliestOccurrence,
		"test":       Test,
	}

	res, err := api.Post[types.#param("int", "actionId", "The action id of the scheduled action")](client, "system/config", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleApplyConfigChannel: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
