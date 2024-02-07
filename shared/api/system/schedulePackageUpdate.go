package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule full package update for several systems.
func SchedulePackageUpdate(cnxDetails *api.ConnectionDetails, Sids []int, EarliestOccurrence $date) (*types.#param("int", "actionId"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sids":       Sids,
		"earliestOccurrence":       EarliestOccurrence,
	}

	res, err := api.Post[types.#param("int", "actionId")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute schedulePackageUpdate: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
