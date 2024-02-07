package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedules an action to apply errata updates to multiple systems.
func ScheduleApplyErrata(cnxDetails *api.ConnectionDetails, Sids []int, ErrataIds []int, AllowModules bool, EarliestOccurrence $date, OnlyRelevant bool, Sid int, OnlyRelevant bool) (*types.#array_single("int", "actionId"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sids":       Sids,
		"errataIds":       ErrataIds,
		"allowModules":       AllowModules,
		"earliestOccurrence":       EarliestOccurrence,
		"onlyRelevant":       OnlyRelevant,
		"sid":       Sid,
		"onlyRelevant":       OnlyRelevant,
	}

	res, err := api.Post[types.#array_single("int", "actionId")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleApplyErrata: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
