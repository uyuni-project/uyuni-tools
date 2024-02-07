package snapshot

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List snapshots for a given system.
 A user may optionally provide a start and end date to narrow the snapshots that
 will be listed.  For example,
 
 If the user provides startDate only, all snapshots created either on or after
 the date provided will be returned.
 If user provides startDate and endDate, all snapshots created on or between the
 dates provided will be returned.
 If the user doesn't provide a startDate and endDate, all snapshots associated
 with the server will be returned.
 
func ListSnapshots(cnxDetails *api.ConnectionDetails, Sid int, StartDate $type, EndDate $type) (*types.#return_array_begin()
      $ServerSnapshotSerializer
  #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system/provisioning/snapshot"
	params := ""
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if StartDate {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if EndDate {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
      $ServerSnapshotSerializer
  #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listSnapshots: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
