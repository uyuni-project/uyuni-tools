package systemgroup

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedules an action to apply errata updates to active systems
 from a group.
func ScheduleApplyErrataToActive(cnxDetails *api.ConnectionDetails, SystemGroupName string, ErrataIds []int, EarliestOccurrence $date, OnlyRelevant bool) (*types.#array_single("int", "actionId"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"systemGroupName":       SystemGroupName,
		"errataIds":       ErrataIds,
		"earliestOccurrence":       EarliestOccurrence,
		"onlyRelevant":       OnlyRelevant,
	}

	res, err := api.Post[types.#array_single("int", "actionId")](client, "systemgroup", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleApplyErrataToActive: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
