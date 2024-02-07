package snapshot

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Deletes all snapshots across multiple systems based on the given date
 criteria.  For example,
 
 If the user provides startDate only, all snapshots created either on or after
 the date provided will be removed.
 If user provides startDate and endDate, all snapshots created on or between the
 dates provided will be removed.
 If the user doesn't provide a startDate and endDate, all snapshots will be
 removed.
 
func DeleteSnapshots(cnxDetails *api.ConnectionDetails, StartDate $type, EndDate $type, Sid int, Sid int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"startDate":       StartDate,
		"endDate":       EndDate,
		"sid":       Sid,
		"sid":       Sid,
	}

	res, err := api.Post[types.#return_int_success()](client, "system/provisioning/snapshot", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteSnapshots: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
